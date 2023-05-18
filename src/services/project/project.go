package project

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"
)

type ProjectType string

const (
	NodeProject ProjectType = "node"
	GoProject   ProjectType = "go"
)

type Project struct {
	Name           string
	Version        string
	Directory      string
	RuntimeVersion string
	Resources      fs.FS
}

func New(name string, directory string, runtimeVersion string, version string, resources fs.FS) Project {
	return Project{
		Name:           name,
		Directory:      directory,
		RuntimeVersion: runtimeVersion,
		Resources:      resources,
		Version:        version,
	}
}

func (proj Project) detectType() (ProjectType, error) {
	if _, err := os.Stat(path.Join(proj.Directory, "package.json")); err == nil {
		return NodeProject, nil
	} else if _, err := os.Stat(path.Join(proj.Directory, "go.mod")); err == nil {
		return GoProject, nil
	} else {
		return "", fmt.Errorf("Cannot detect project type %v", err)
	}
}

func (proj Project) loadDockerfile() ([]byte, error) {
	projectType, err := proj.detectType()
	if err != nil {
		return []byte{}, err
	}

	dockerfileTemplate := path.Join("assets", "docker", fmt.Sprintf("Dockerfile.%s.tmpl", projectType))
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
func (proj Project) Dockerfile() ([]byte, error) {
	content, err := proj.loadDockerfile()
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}
