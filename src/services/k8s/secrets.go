package k8s

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/nearform/initium-cli/src/utils/logger"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func ListSecrets(config *rest.Config, namespace string) (*v1.SecretList, error) {
	convertSecretListToString := func (secretList *v1.SecretList) string {
	    var sb strings.Builder

	    for _, secret := range secretList.Items {
	        sb.WriteString(fmt.Sprintf("%s\n", secret.ObjectMeta.Name))
	        for key := range secret.Data {
	            sb.WriteString(fmt.Sprintf("  %s\n", key))
	        }
	        sb.WriteString("\n")
	    }

	    return sb.String()
	}

	ctx := context.Background()
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Creating K8s client %v", err)
	}

	secretList, err := client.CoreV1().Secrets(namespace).List(ctx, metav1.ListOptions{})

	if err != nil {
		return nil, fmt.Errorf("Error fetching K8s secret %v", err.Error())
	}

	stringSecretList := convertSecretListToString(secretList);
	logger.PrintInfo("*********** CODEIUM CONVERTED"); // TODO: Return this and print it from CMD
	logger.PrintInfo(stringSecretList); // TODO: Return this and print it from CMD
	logger.PrintInfo("*********** SECRETLIST.STRING NATIVE METHOD"); // TODO: Return this and print it from CMD
	logger.PrintInfo(secretList.String()); // TODO: Return this and print it from CMD

	return secretList, nil
}

func GetSecret(config *rest.Config, secretName string, namespace string) (*v1.Secret, error) {
	ctx := context.Background()
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Creating K8s client %v", err)
	}

	secret, err := getSecretUsingExistingCtx(client, ctx, secretName, namespace)
	if err == nil {
		return nil, fmt.Errorf("Error fetching K8s secret %v", err)
	}
	return secret, nil
}

func UpdateSecretKeyValue(config *rest.Config, secretName string, secretKey string, secretValue string, namespace string) error {
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
