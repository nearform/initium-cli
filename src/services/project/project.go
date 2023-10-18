package project

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"text/template"

	"github.com/nearform/initium-cli/src/services/git"
	"github.com/nearform/initium-cli/src/utils/defaults"
)

type ProjectType string

const (
	NodeProject ProjectType = "node"
	GoProject   ProjectType = "go"
)

type Project struct {
	Name                  string
	Version               string
	Type                  ProjectType
	Directory             string
	RuntimeVersion        string
	DefaultRuntimeVersion string
	ImagePullSecrets      []string
	Resources             fs.FS
}

type InitOptions struct {
	DestinationFolder string
	DefaultBranch     string
	PipelineType      string
}

func GuessAppName() *string {
	var name string
	name, err := git.GetRepoName()
	if err != nil {
		return nil
	}
	return &name
}

func New(name string, projectType ProjectType, directory string, runtimeVersion string, version string, imagePullSecrets []string, resources fs.FS) Project {
	return Project{
		Name:             name,
		Type:             projectType,
		Directory:        directory,
		RuntimeVersion:   runtimeVersion,
		ImagePullSecrets: imagePullSecrets,
		Resources:        resources,
		Version:          version,
	}
}

func DetectType(directory string) (ProjectType, error) {
	var detectedRuntimes []ProjectType
	var projectType ProjectType
	if _, err := os.Stat(path.Join(directory, "package.json")); err == nil {
		detectedRuntimes = append(detectedRuntimes, NodeProject)
		projectType = NodeProject
	}
	if _, err := os.Stat(path.Join(directory, "go.mod")); err == nil {
		detectedRuntimes = append(detectedRuntimes, GoProject)
		projectType = GoProject
	}
	if len(detectedRuntimes) == 0 {
		return "", fmt.Errorf("cannot detect project type")
	}
	if len(detectedRuntimes) > 1 {
		return "", fmt.Errorf("more than one project runtimes detected (%v), use --project-type flag or the INITIUM_PROJECT_TYPE env variable to set the desired runtime", detectedRuntimes)
	}
	return projectType, nil
}

func (proj *Project) matchType() (ProjectType, error) {
	switch proj.Type {
	case NodeProject:
		return NodeProject, nil
	case GoProject:
		return GoProject, nil
	default:
		return "", fmt.Errorf("cannot detect project type %s", proj.Type)
	}
}

func (proj *Project) setRuntimeVersion() error {
	switch proj.Type {
	case NodeProject:
		proj.DefaultRuntimeVersion = defaults.DefaultNodeRuntimeVersion
		return nil
	case GoProject:
		proj.DefaultRuntimeVersion = defaults.DefaultGoRuntimeVersion
		return nil
	default:
		return fmt.Errorf("cannot detect runtime version for project type %s", proj.Type)
	}
}

func (proj Project) loadDockerfile() ([]byte, error) {
	var projectType ProjectType
	var err error
	if proj.Type != "" {
		projectType, err = proj.matchType()
	} else {
		projectType, err = DetectType(proj.Directory)
		proj.Type = projectType
	}
	if err != nil {
		return []byte{}, err
	}
	err = proj.setRuntimeVersion()
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

		destinationFile := path.Join(options.DestinationFolder, fmt.Sprintf("initium_%s.yaml", tmpl))

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
