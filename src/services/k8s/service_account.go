package k8s

import (
	"fmt"
	"io/fs"
	"path"
)

func GetServiceAccount(resources fs.FS) error {
	serviceAccountManifest, err := getServiceAccountManifest(resources)
	if err != nil {
		return err
	}

	fmt.Printf("%s", serviceAccountManifest)

	return nil
}

func getServiceAccountManifest(resources fs.FS) (string, error) {
	manifest := ""
	for _, v := range []string{"cluster-role", "role-binding", "service-account", "token"} {
		fileContent, err := getFileContent(v, resources)
		if err != nil {
			return "", err
		}
		manifest = manifest + fmt.Sprintf("---\n%s\n", fileContent)
	}

	return manifest, nil
}

func getFileContent(filename string, resources fs.FS) ([]byte, error) {
	k8sTemplatesPath := path.Join("assets", "k8s", "serviceAccount")

	return fs.ReadFile(resources, path.Join(k8sTemplatesPath, filename+".yaml"))
}
