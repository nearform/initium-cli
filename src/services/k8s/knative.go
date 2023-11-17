package k8s

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/services/project"
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

func LoadManifest(namespace string, commitSha string, project *project.Project, dockerImage docker.DockerImage, envFile string, secretRefEnvFile string) (*servingv1.Service, error) {
	manifestEnvVars := map[string]string{}
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
	if err = setEnv(service, envFile, manifestEnvVars); err != nil {
		return nil, err
	}
	if err = setSecretEnv(service, secretRefEnvFile, manifestEnvVars); err != nil {
		return nil, err
	}

	return service, nil
}

func setSecretEnv(manifest *servingv1.Service, secretRefEnvFile string, manifestEnvVars map[string]string) error {
	secretEnvVarList, err := loadEnvFile(secretRefEnvFile, manifestEnvVars)
	if err != nil {
		return err
	}
	for _, secretEnvVar := range secretEnvVarList { //eg: [MOCK5=kubernetessecretname/secretkey]
		err := validateSecretEnvVar(secretEnvVar)
		if err != nil {
			return err
		}

		secretValue := strings.SplitN(secretEnvVar.Value, "/", 2)
		manifest.Spec.Template.Spec.Containers[0].Env = append(manifest.Spec.Template.Spec.Containers[0].Env, corev1.EnvVar{
			Name: secretEnvVar.Name, //eg: MOCK5
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					Key: secretValue[1], //eg: kubernetesecretkey
					LocalObjectReference: corev1.LocalObjectReference{
						Name: secretValue[0], //eg: kubernetesecretname
					},
				},
			},
		})
	}
	return nil
}

func validateSecretEnvVar(secretEnvVar corev1.EnvVar) error {
	// Mandatory char
	if !strings.Contains(secretEnvVar.Value, "/") {
		return fmt.Errorf("Invalid secret format for '%s'. Missing '/' char. Value must be in the format <secret-name>/<secret-key>, instead of '%s'", secretEnvVar.Name, secretEnvVar.Value)
	}
	return nil
}

func setEnv(manifest *servingv1.Service, envFile string, manifestEnvVars map[string]string) error {
	envVarList, err := loadEnvFile(envFile, manifestEnvVars)
	if err != nil {
		return err
	}
	manifest.Spec.Template.Spec.Containers[0].Env = append(manifest.Spec.Template.Spec.Containers[0].Env, envVarList...)
	return nil
}

func loadEnvFile(envFile string, manifestEnvVars map[string]string) ([]corev1.EnvVar, error) {
	var envVarList []corev1.EnvVar
	
	if _, err := os.Stat(envFile); errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	
	envVariables, err := godotenv.Read(envFile)
	if err != nil {
		return nil, fmt.Errorf("Error loading .env file. '%s' already set", err)
	}

	if len(envVariables) > 0 {
		for key, value := range envVariables {
			if manifestEnvVars[key] == "" {
				manifestEnvVars[key] = value
				envVar := corev1.EnvVar{
					Name:  key,
					Value: value,
				}
				envVarList = append(envVarList, envVar)
			} else {
				return nil, fmt.Errorf("Conflicting environment variable. '%s' already set through another file", key)
			}
		}
		log.Infof("Environment variables file %v content is now loaded!", envFile)
	} else {
		log.Warnf("Environment file %v is empty, Nothing to load!", envFile)
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
			return fmt.Errorf("%s", service.Status.URL)
		}

		time.Sleep(time.Millisecond * 500)
	}
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
