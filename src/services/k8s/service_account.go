package k8s

import (
	"fmt"
	"io/fs"
	"path"
)

func GetServiceAccount(resources fs.FS) error {
	k8sTemplatesPath := path.Join("assets", "k8s", "serviceAccount")

	for _, v := range []string{"cluster-role", "role-binding", "service-account", "token"} {
		filecontent, err := fs.ReadFile(resources, path.Join(k8sTemplatesPath, v+".yaml"))
		if err != nil {
			return err
		}

		fmt.Println("---")
		fmt.Printf("%s\n", filecontent)
	}

	return nil
}
