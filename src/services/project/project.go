package project

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
)

type ProjectType string

const (
	NodeProject ProjectType = "node"
	GoProject   ProjectType = "go"
)

type Project struct {
	Name                  string
	Version               string
	Directory             string
	RuntimeVersion        string
	DefaultRuntimeVersion string
	Resources             fs.FS
}

type InitOptions struct {
	DestinationFolder string
	DefaultBranch     string
	PipelineType      string
	Repository        string
	AppName           string
	ProjectDirectory  string
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

func (proj *Project) detectType() (ProjectType, error) {
	if _, err := os.Stat(path.Join(proj.Directory, "package.json")); err == nil {
		proj.DefaultRuntimeVersion = defaults.DefaultNodeRuntimeVersion
		return NodeProject, nil
	} else if _, err := os.Stat(path.Join(proj.Directory, "go.mod")); err == nil {
		proj.DefaultRuntimeVersion = defaults.DefaultGoRuntimeVersion
		return GoProject, nil
	} else {
		return "", fmt.Errorf("cannot detect project type %v", err)
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
	if err = template.Execute(output, proj); err != nil {
		return []byte{}, err
	}
	return output.Bytes(), nil
}

func (proj Project) Dockerfile() ([]byte, error) {
	content, err := proj.loadDockerfile()
	if err != nil {
		return []byte{}, err
	}

	return content, nil
}

func ProjectInit(options InitOptions, resources fs.FS) ([]string, error) {

	returnData := []string{}
	for _, tmpl := range []string{"onmain", "onbranch"} {
		template, err := template.ParseFS(resources, path.Join("assets", options.PipelineType, fmt.Sprintf("%s.tmpl", tmpl)))

		if err != nil {
			return returnData, fmt.Errorf("error: %v", err)
		}

		fileContent := &bytes.Buffer{}
		if err = template.Execute(fileContent, options); err != nil {
			return returnData, err
		}

		destinationFile := path.Join(options.DestinationFolder, fmt.Sprintf("kka_%s.yaml", tmpl))

		if err := os.MkdirAll(options.DestinationFolder, os.ModePerm); err != nil {
			return returnData, fmt.Errorf("error: %v", err)
		}

		// I assume that the file is in source control and the user will be able to
		// revert the changes, I'll create an issue to make this step interactive so
		// we can ask confirmation to override the file.
		if err = os.WriteFile(destinationFile, fileContent.Bytes(), 0755); err != nil {
			return returnData, fmt.Errorf("error: %v", err)
		}

		returnData = append(returnData, destinationFile)
	}

	return returnData, nil
}

func (proj Project) NodeInstallCommand() string {
	installCommand := "npm i"

	if _, err := os.Stat(path.Join(proj.Directory, "package-lock.json")); err == nil {
		installCommand = "npm ci"
	}
	return installCommand
}
