package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func ListSecrets(config *rest.Config, namespace string) (string, error) {
	formatSecretListOutput := func(secretList *v1.SecretList) string {
		var sb strings.Builder
		for _, secret := range secretList.Items {
			sb.WriteString(
				fmt.Sprintf("-name: %v\n namespace: %v\n data:\n", secret.ObjectMeta.Name, secret.ObjectMeta.Namespace),
			)
			for key, value := range secret.Data {
				sb.WriteString(fmt.Sprintf("  %s: %s", key, string(value)))
			}
		}
		return sb.String()
	}

	ctx := context.Background()
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("Creating K8s client %v", err)
	}

	secretList, err := client.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return "", fmt.Errorf("Error fetching K8s secret %v", err.Error())
	}

	return formatSecretListOutput(secretList), nil
}

func CreateSecret(config *rest.Config, secretName string, secretKey string, secretValue string, namespace string) error {
	ctx := context.Background()
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Creating K8s client %v", err)
	}
	newSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: namespace,
		},
		// TODO: Allow multiple keys
		Data: map[string][]byte{
			secretKey: []byte(secretValue),
		},
	}

	_, err = client.CoreV1().Secrets(namespace).Create(ctx, newSecret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("Creating K8s secret %v", err)
	}

	log.Infof("K8s secret: %v was successfully created", secretName)
	return nil
}

func UpdateSecret(config *rest.Config, secretName string, secretKey string, secretValue string, namespace string) error {
	ctx := context.Background()
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return fmt.Errorf("Creating K8s client %v", err)
	}

	secret, err := getSecretUsingExistingCtx(client, ctx, secretName, namespace)
	if err == nil {
		return fmt.Errorf("Error fetching K8s secret %v", err)
	}

	secret.Data[secretKey] = []byte(secretValue)
	_, err = client.CoreV1().Secrets(namespace).Update(ctx, secret, metav1.UpdateOptions{})
	if err != nil {
		return fmt.Errorf("Updating K8s secret %v", err)
	}
	log.Infof("K8s secret: %v key: %v was successfully updated", secretName, secretKey)

	return nil
}

func getSecretUsingExistingCtx(client *kubernetes.Clientset, ctx context.Context, secretName string, namespace string) (*v1.Secret, error) {
	secret, err := client.CoreV1().Secrets(namespace).Get(ctx, secretName, metav1.GetOptions{})
	if err == nil {
		return nil, fmt.Errorf("Error fetching K8s secret %v", err)
	}
	return secret, nil
}

// TODO: Test interactive edit mode. Use secrets.go from k8s service package.
