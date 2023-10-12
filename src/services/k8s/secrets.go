package k8s

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func SetSecret(secretName string, secretKey string, secretValue string, config *rest.Config, namespace string) error {

	log.Info("Creating K8s secret..", "secret name", secretName)
	ctx := context.Background()

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Creating K8s client %v", err)
	}

	secret, err := client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err != nil {
		newSecret := &corev1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      secretName,
				Namespace: namespace,
			},
			Data: map[string][]byte{
				secretKey: []byte(secretValue),
			},
		}
		_, err = client.CoreV1().Secrets(namespace).Create(ctx, newSecret, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("Creating K8s secret %v", err)
		}
		log.Infof("K8s secret: %v was successfully created", secretName)
	} else {
		secret.Data = map[string][]byte{
			secretKey: []byte(secretValue),
		}

		_, err = client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
		if err != nil {
			return fmt.Errorf("Updating K8s secret %v", err)
		}
		log.Infof("K8s secret: %v was successfully updated", secretName)
	}
	return nil
}
