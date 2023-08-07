package k8s

import (
	"embed"
	"fmt"
	"path"
)

func GetServiceAccount(resources embed.FS) error {
	k8sTemplatesPath := path.Join("assets", "k8s", "serviceAccount")

	for _, v := range []string{"cluster-role", "role-binding", "service-account", "token"} {
		filecontent, err := resources.ReadFile(path.Join(k8sTemplatesPath, v+".yaml"))
		if err != nil {
			return err
		}

		fmt.Println("---")
		fmt.Printf("%s\n", filecontent)
	}

	return nil
}
