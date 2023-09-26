package k8s

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"text/template"
	"time"
    "strings"
	"bufio"
	"os"
	
	"github.com/charmbracelet/log"
	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/services/project"

	corev1 "k8s.io/api/core/v1"
	apimachineryErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	certutil "k8s.io/client-go/util/cert"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	servingv1client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

func Config(endpoint string, token string, caCrt []byte) (*rest.Config, error) {
	if _, err := certutil.NewPoolFromBytes(caCrt); err != nil {
		return nil, fmt.Errorf("Expected to load root CA from bytes, but got err: %v", err)
	}
	return &rest.Config{
		// TODO: switch to using cluster DNS.
		Host: endpoint,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: caCrt,
		},
		BearerToken: string(token),
	}, nil
}

func loadManifest(project *project.Project, dockerImage docker.DockerImage) (*servingv1.Service, error) {
	knativeTemplate := path.Join("assets", "knative", "service.yaml.tmpl")
	template, err := template.ParseFS(project.Resources, knativeTemplate)
	if err != nil {
		return nil, fmt.Errorf("error reading the knative service yaml: %v", err)
	}

	output := &bytes.Buffer{}
	// TODO replace map[string]string{} with proper values
	if err = template.Execute(output, dockerImage); err != nil {
		return nil, err
	}

	data := output.Bytes()

	err = servingv1.AddToScheme(scheme.Scheme)
	if err != nil {
		return nil, fmt.Errorf("error adding Knative Serving scheme: %v", err)
	}

	// Decode the YAML data into a Knative Service object
	obj, _, err := scheme.Codecs.UniversalDeserializer().Decode(data, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("error decoding YAML: %v", err)
	}

	service, ok := obj.(*servingv1.Service)
	if !ok {
		return nil, fmt.Errorf("decoded object is not a Knative Service: %v", obj)
	}

	return service, nil
}

func Apply(namespace string, config *rest.Config, project *project.Project, dockerImage docker.DockerImage) error {
	log.Info("Deploying Knative service", "host", config.Host, "name", project.Name, "namespace", namespace)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	serviceManifest, err := loadManifest(project, dockerImage)
	if err != nil {
		return err
	}

	// Create a new Knative Serving client
	servingClient, err := servingv1client.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Error creating the knative client %v", err)
	}

	serviceManifest.ObjectMeta.Namespace = namespace
	serviceManifest.ObjectMeta.Name = project.Name

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Creating Kubernetes client %v", err)
	}

	_, err = client.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceManifest.ObjectMeta.Namespace,
		},
	}, metav1.CreateOptions{})

	if err != nil && !apimachineryErrors.IsAlreadyExists(err) {
		return fmt.Errorf("cannot create namespace %s, failed with %v", serviceManifest.ObjectMeta.Namespace, err)
	}

	if _, err := os.Stat(".env"); err != nil {
        if os.IsNotExist(err) {
			log.Info("No environment variables file (.env) to Load!")
        } else {
            log.Fatalf("Error checking .env file: %v", err)
        } 
	} else {
		log.Info("Environment variables file (.env) found! Loading..")
		file, err := os.Open(".env")
		if err != nil {
			log.Fatalf("Error opening .env file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		envVariables := make(map[string]string)
		
		checkFormat := func(line string) bool {
			parts := strings.SplitN(line, "=", 2)
			return len(parts) == 2
		}

		for scanner.Scan() {
			line := scanner.Text()
				if checkFormat(line) {
				parts := strings.SplitN(line, "=", 2)
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
					envVariables[key] = value
			} else {
				log.Warnf("Environment variables file (.env) line won't be processed due to invalid format: %s. Accepted: KEY=value", line)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Errorf("Error reading environment variables file (.env): %v", err)
		}

		if len(envVariables) > 0 {
			var envVarList []corev1.EnvVar
			for key, value := range envVariables {
				envVar := corev1.EnvVar{
					Name:  key,
					Value: value,
				}
				envVarList = append(envVarList, envVar)
			}
			serviceManifest.Spec.Template.Spec.Containers[0].Env = append(serviceManifest.Spec.Template.Spec.Containers[0].Env, envVarList...)
			log.Info("Environment variables file (.env) content is now loaded!")
		} else {
			log.Warnf("Environment file (.env) is empty, Nothing to load!")
		}
	}
	
	service, err := servingClient.Services(serviceManifest.ObjectMeta.Namespace).Get(ctx, serviceManifest.ObjectMeta.Name, metav1.GetOptions{})
	var deployedService *servingv1.Service
	if err != nil {
		deployedService, err = servingClient.Services(serviceManifest.ObjectMeta.Namespace).Create(ctx, serviceManifest, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Creating Knative service %v", err)
		}
	} else {
		service.Spec = serviceManifest.Spec
		deployedService, err = servingClient.Services(serviceManifest.ObjectMeta.Namespace).Update(ctx, service, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("Updating Knative service %v", err)
		}
	}

	fmt.Printf("Knative service %q deployed successfully in namespace %s.\n", deployedService.GetObjectMeta().GetName(), serviceManifest.ObjectMeta.Namespace)

	for {
		service, err = servingClient.Services(serviceManifest.ObjectMeta.Namespace).Get(ctx, serviceManifest.ObjectMeta.Name, metav1.GetOptions{})
		if err != nil {
			return err
		}
		if service.Status.URL != nil {
			fmt.Printf("You can reach it via %s\n", service.Status.URL)
			break
		}

		time.Sleep(time.Millisecond * 500)
	}

	return nil
}

func Clean(namespace string, config *rest.Config, project *project.Project) error {
	log.Info("Deleting Knative service", "host", config.Host, "name", project.Name, "namespace", namespace)
	ctx := context.Background()

	// Create a new Knative Serving client
	servingClient, err := servingv1client.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Error creating the knative client %v", err)
	}

	err = servingClient.Services(namespace).Delete(ctx, project.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deleting the service: %v", err)
	}

	log.Info("The Knative service was successfully deleted", "host", config.Host, "name", project.Name, "namespace", namespace)
	return nil
}
