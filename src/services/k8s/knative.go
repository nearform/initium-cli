package k8s

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/services/project"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"sigs.k8s.io/yaml"

	corev1 "k8s.io/api/core/v1"
	apimachineryErrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	certutil "k8s.io/client-go/util/cert"
	servingv1 "knative.dev/serving/pkg/apis/serving/v1"
	servingv1client "knative.dev/serving/pkg/client/clientset/versioned/typed/serving/v1"
)

const (
	UpdateShaAnnotationName       = "initium.nearform.com/updateSha"
	UpdateTimestampAnnotationName = "initium.nearform.com/updateTimestamp"
	visibilityLabel               = "networking.knative.dev/visibility"
	visibilityLabelPrivateValue   = "cluster-local"
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

func setLabels(manifest *servingv1.Service, project project.Project) {
	if manifest.ObjectMeta.Labels == nil {
		manifest.ObjectMeta.Labels = map[string]string{}
	}
	if project.IsPrivate {
		manifest.ObjectMeta.Labels[visibilityLabel] = visibilityLabelPrivateValue
	}
}

func LoadManifest(namespace string, commitSha string, project *project.Project, dockerImage docker.DockerImage, envFile string) (*servingv1.Service, error) {
	knativeTemplate := path.Join("assets", "knative", "service.yaml.tmpl")
	template, err := template.ParseFS(project.Resources, knativeTemplate)
	if err != nil {
		return nil, fmt.Errorf("error reading the knative service yaml: %v", err)
	}

	templateParams := map[string]interface{}{
		"Name":             dockerImage.Name,
		"RemoteTag":        dockerImage.RemoteTag(),
		"ImagePullSecrets": project.ImagePullSecrets,
	}

	output := &bytes.Buffer{}
	// TODO replace map[string]string{} with proper values
	if err = template.Execute(output, templateParams); err != nil {
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

	service.ObjectMeta.Namespace = namespace
	service.ObjectMeta.Name = project.Name
	service.Spec.Template.ObjectMeta.Annotations = map[string]string{
		UpdateShaAnnotationName:       commitSha,
		UpdateTimestampAnnotationName: time.Now().Format(time.RFC3339),
	}

	setLabels(service, *project)
	if err = setEnv(service, envFile); err != nil {
		return nil, err
	}

	return service, nil
}

func setEnv(manifest *servingv1.Service, envFile string) error {
	envVarList, err := loadEnvFile(envFile)
	if err != nil {
		return err
	}

	manifest.Spec.Template.Spec.Containers[0].Env = append(manifest.Spec.Template.Spec.Containers[0].Env, envVarList...)
	return nil
}

func loadEnvFile(envFile string) ([]corev1.EnvVar, error) {
	var envVarList []corev1.EnvVar
	if _, err := os.Stat(envFile); err != nil {
		if (os.IsNotExist(err)) && (path.Base(envFile) == defaults.EnvVarFile) {
			log.Infof("No environment variables file %s to Load!", defaults.EnvVarFile)
		} else {
			return nil, fmt.Errorf("Error loading %v file: %v", envFile, err)
		}
	} else {
		log.Infof("Environment variables file %s found! Loading..", envFile)
		file, err := os.Open(envFile)
		if err != nil {
			return nil, fmt.Errorf("Error opening %v file: %v", envFile, err)
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
				log.Warnf("Environment variables file %v line won't be processed due to invalid format: %s. Accepted: KEY=value", envFile, line)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("Error reading environment variables file %v: %v", envFile, err)
		}

		if len(envVariables) > 0 {
			for key, value := range envVariables {
				envVar := corev1.EnvVar{
					Name:  key,
					Value: value,
				}
				envVarList = append(envVarList, envVar)
			}
			log.Infof("Environment variables file %v content is now loaded!", envFile)
		} else {
			log.Warnf("Environment file %v is empty, Nothing to load!", envFile)
		}
	}
	return envVarList, nil
}

func ToYaml(serviceManifest *servingv1.Service) ([]byte, error) {
	scheme := runtime.NewScheme()
	servingv1.AddToScheme(scheme)
	codec := serializer.NewCodecFactory(scheme).LegacyCodec(servingv1.SchemeGroupVersion)
	jsonBytes, err := runtime.Encode(codec, serviceManifest)
	if err != nil {
		return nil, err
	}
	return yaml.JSONToYAML(jsonBytes)
}

func Apply(serviceManifest *servingv1.Service, config *rest.Config) error {
	log.Info("Deploying Knative service", "host", config.Host, "name", serviceManifest.ObjectMeta.Name, "namespace", serviceManifest.ObjectMeta.Namespace)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// Create a new Knative Serving client
	servingClient, err := servingv1client.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Error creating the knative client %v", err)
	}

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
			exec.Command(fmt.Sprintf("echo \"INITIUM_OUTPUT_URL=%s\" >> \"$GITHUB_OUTPUT\"", service.Status.URL))
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
