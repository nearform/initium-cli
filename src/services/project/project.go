package project

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"text/template"
)

type Project struct {
	Name           string
	Version        string
	Directory      string
	RuntimeVersion string
	Resources      embed.FS
}

func New(name string, directory string, runtimeVersion string, version string, resources embed.FS) Project {
	return Project{
		Name:           name,
		Directory:      directory,
		RuntimeVersion: runtimeVersion,
		Resources:      resources,
		Version:        version,
	}
}

func (proj Project) detectType() (string, error) {
	if _, err := os.Stat(path.Join(proj.Directory, "package.json")); err == nil {
		return "node", nil
	} else if _, err := os.Stat(path.Join(proj.Directory, "go.mod")); err == nil {
		return "go", nil
	} else {
		return "", fmt.Errorf("Cannot detect project type %v", err)
	}
}

func (proj Project) loadDockerfile() ([]byte, error) {
	projectType, err := proj.detectType()
	if err != nil {
		return []byte{}, err
	}

	dockerfileTemplate := path.Join("assets", "docker", "Dockerfile."+projectType+".tmpl")
	template, err := template.ParseFS(proj.Resources, dockerfileTemplate)
	if err != nil {
		return []byte{}, err
	}

	output := &bytes.Buffer{}
	// TODO replace map[string]string{} with proper values
	if err = template.Execute(output, proj); err != nil {
		return []byte{}, err
	}
	return output.Bytes(), nil
}

// TODO: there is no need to persist this file, we could add it to the tar context from memory or a temp dir
func (proj Project) AddDockerFile() error {
	content, err := proj.loadDockerfile()
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
