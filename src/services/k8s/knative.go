package k8s

import (
	"context"
	"fmt"
	"path"
	"time"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func loadManifest(project project.Project) (*servingv1.Service, error) {
	data, err := project.Resources.ReadFile(path.Join("assets", "knative", "service.yaml"))
	if err != nil {
		return nil, fmt.Errorf("error reading the knative service yaml: %v", err)
	}

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

func Apply(config *rest.Config, project project.Project) error {
	logger.PrintInfo("Deploying Knative service to " + config.Host)
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()
	// Create a new Knative Serving client
	servingClient, err := servingv1client.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Error creating the knative client %v", err)
	}

	service, err := loadManifest(project)
	if err != nil {
		return err
	}

	service.ObjectMeta.Namespace = project.Version
	service.ObjectMeta.Name = project.Name

	getService, err := servingClient.Services(service.ObjectMeta.Namespace).Get(ctx, service.ObjectMeta.Name, metav1.GetOptions{})

	if err != nil {
		createdService, err := servingClient.Services(service.ObjectMeta.Namespace).Create(ctx, service, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Creating Knative service %v", err)
		}
		fmt.Printf("Created Knative service %q.\n", createdService.GetObjectMeta().GetName())
	} else {
		updatedService, err := servingClient.Services(service.ObjectMeta.Namespace).Update(ctx, getService, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("Updating Knative service %v", err)
		}
		fmt.Printf("Updated Knative service %q.\n", updatedService.GetObjectMeta().GetName())
	}

	return nil
}

func Clean(config *rest.Config, project project.Project) error {
	logger.PrintInfo("Deleting Knative service from " + config.Host)
	ctx := context.Background()

	// Create a new Knative Serving client
	servingClient, err := servingv1client.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Error creating the knative client %v", err)
	}

	err = servingClient.Services(project.Version).Delete(ctx, project.Name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("deleting the service: %v", err)
	}

	logger.PrintInfo("Deleted Knative service " + project.Name)
	return nil
}
