package project

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Project struct {
	Name      string
	Directory string
}

func New(projectName string, projectDirectory string) Project {
	return Project{
		Name:      projectName,
		Directory: projectDirectory,
	}
}

func (proj Project) detectType() (string, error) {
	if _, err := os.Stat(path.Join(proj.Directory, "package.json")); err == nil {
		return "node", nil
	} else {
		return "", fmt.Errorf("Cannot detect project type %v", err)
	}
}

func (proj Project) loadDockerfile(resources embed.FS) ([]byte, error) {
	projectType, err := proj.detectType()
	if err != nil {
		return []byte{}, err
	}
	dockerfileTemplate := path.Join(".", "dockerfiles", "Dockerfile."+projectType)
	return resources.ReadFile(dockerfileTemplate)
}

// TODO: there is no need to persist this file, we could add it to the tar context from memory or a temp dir
func (proj Project) AddDockerFile(resources embed.FS) error {
	content, err := proj.loadDockerfile(resources)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(path.Join(proj.Directory, "Dockerfile.kka"), content, 0644)
	if err != nil {
		return fmt.Errorf("Writing Dockerfile content: %v", err)
	}
	return nil
}

func (proj Project) DeleteDockerFile() error {
	err := os.Remove(path.Join(proj.Directory, "Dockerfile.kka"))
	if err != nil {
		return fmt.Errorf("Deleting generated dockerfile: %v", err)
	}
	return nil
}
