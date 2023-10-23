package project

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/nearform/initium-cli/src/services/git"
	"github.com/nearform/initium-cli/src/utils/defaults"
)

type ProjectType string

const (
	NodeProject       ProjectType = "node"
	GoProject         ProjectType = "go"
	FrontendJsProject ProjectType = "frontend-js"
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
	IsPrivate             bool
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

func New(name string, projectType ProjectType, directory string, runtimeVersion string, version string, isPrivate bool, imagePullSecrets []string, resources fs.FS) Project {
	return Project{
		Name:             name,
		Type:             projectType,
		Directory:        directory,
		RuntimeVersion:   runtimeVersion,
		ImagePullSecrets: imagePullSecrets,
		Resources:        resources,
		Version:          version,
		IsPrivate:        isPrivate,
	}
}

func IsValidProjectType(projectType string) bool {
	switch projectType {
	case string(NodeProject):
		return true
	case string(GoProject):
		return true
	case string(FrontendJsProject):
		return true
	default:
		return false
	}
}

func DetectType(directory string) (ProjectType, error) {
	var detectedRuntimes []ProjectType
	var projectType ProjectType
	if _, err := os.Stat(path.Join(directory, "package.json")); err == nil {
		bytes, err := os.ReadFile(path.Join(directory, "package.json"))
		if err != nil {
			fmt.Print(err)
		}
		fileStr := string(bytes)
		if strings.Contains(fileStr, "react") || strings.Contains(fileStr, "angular") || strings.Contains(fileStr, "vue") {
			detectedRuntimes = append(detectedRuntimes, FrontendJsProject)
			projectType = FrontendJsProject
		} else {
			detectedRuntimes = append(detectedRuntimes, NodeProject)
			projectType = NodeProject
		}
	}
	if _, err := os.Stat(path.Join(directory, "go.mod")); err == nil {
		detectedRuntimes = append(detectedRuntimes, GoProject)
		projectType = GoProject
	}
	if len(detectedRuntimes) == 0 {
		return "", fmt.Errorf("cannot detect the project type by checking the repository file structure, use the --project-type flag or the INITIUM_PROJECT_TYPE env variable to set the desired runtime")
	}
	if len(detectedRuntimes) > 1 {
		return "", fmt.Errorf("more than one project runtimes detected (%v), use the --project-type flag or the INITIUM_PROJECT_TYPE env variable to set the desired runtime", detectedRuntimes)
	}
	return projectType, nil
}

func (proj *Project) setRuntimeVersion() error {
	switch proj.Type {
	case NodeProject:
		proj.DefaultRuntimeVersion = defaults.DefaultNodeRuntimeVersion
		return nil
	case GoProject:
		proj.DefaultRuntimeVersion = defaults.DefaultGoRuntimeVersion
		return nil
	case FrontendJsProject:
		proj.DefaultRuntimeVersion = defaults.DefaultFrontendJsRuntimeVersion
		return nil
	default:
		return fmt.Errorf("cannot detect runtime version for project type %s", proj.Type)
	}
}

func (proj Project) loadDockerfile() ([]byte, error) {
	var err error
	projectType := proj.Type
	if !IsValidProjectType(string(projectType)) {
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
