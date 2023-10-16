package project

import (
	"fmt"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/nearform/initium-cli/src/utils/defaults"
)

var projects = map[ProjectType]map[string]string{
	NodeProject: {"directory": "example"},
	GoProject:   {"directory": "."},
	"invalid":   {"directory": "src"},
}

var root = "../../../"

func TestDetectType(t *testing.T) {

	for project_type, props := range projects {
		test_proj_type := Project{Name: string(project_type),
			Directory: path.Join(root, props["directory"])}

		var proj_type ProjectType
		var err error
		if test_proj_type.Type == "" {
			proj_type, err = test_proj_type.detectType()
		} else {
			proj_type, err = test_proj_type.matchType()
		}

		// if we cannot autodetect a project we will return an error
		if project_type == "invalid" && err != nil {
			return
		}

		if err != nil {
			t.Fatalf(fmt.Sprintf("Error: %s", err))
		}

		if proj_type != project_type {
			t.Fatalf("Error: %s project not found", project_type)
		}
	}
}

func TestLoadDockerfile(t *testing.T) {
	for project_type, props := range projects {
		proj_dockerfile := Project{Name: string(project_type),
			Directory: path.Join(root, props["directory"]),
			Resources: os.DirFS(root),
		}
		_, err := proj_dockerfile.loadDockerfile()

		// if we cannot autodetect a project we will return an error
		if project_type == "invalid" && err != nil {
			return
		}

		if err != nil {
			t.Fatalf(fmt.Sprintf("Error: %s", err))
		}
	}
}

func TestCorrectRuntime(t *testing.T) {

	proj_runtime := Project{
		Name:           "test",
		Directory:      path.Join(root, projects["node"]["directory"]),
		Resources:      os.DirFS(root),
		RuntimeVersion: "30",
	}
	data, err := proj_runtime.loadDockerfile()

	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	if isContain := strings.Contains(string(data), "node:"+proj_runtime.RuntimeVersion); !isContain {
		t.Fatalf(fmt.Sprintf("Runtime %v not properly replaced", proj_runtime.RuntimeVersion))
	}
}

func TestInit(t *testing.T) {

	supportedPipelineTypes := []string{"github"}

	for _, pipelineType := range supportedPipelineTypes {
		sourcePipelineTemplate := path.Join(root, "assets", pipelineType, "onmain.tmpl")

		// check if embedded folder and pipeline.tmpl file exists
		if _, err := os.Stat(sourcePipelineTemplate); err != nil {
			t.Fatalf(fmt.Sprintf("Error: template file for supported pipeline %v dont exists", pipelineType))
		}

		testDestinationFolder := fmt.Sprintf(".init_test_%s", pipelineType)

		// test function
		files, err := ProjectInit(InitOptions{
			DestinationFolder: testDestinationFolder,
			DefaultBranch:     defaults.GithubDefaultBranch,
			PipelineType:      "github",
		}, os.DirFS(root))

		if err != nil {
			t.Errorf(fmt.Sprintf("Error: %v", err))
		}

		for _, file := range files {
			if _, err := os.Stat(file); os.IsNotExist(err) {
				t.Fatalf(fmt.Sprintf("error: cannot find file %s", file))
			}
			if err = os.Remove(file); err != nil {
				t.Errorf(fmt.Sprintf("error while removing %s: %v", file, err))
			}
		}

		err = os.Remove(testDestinationFolder)
		if err != nil {
			t.Errorf(fmt.Sprintf("error while removing %s: %v", testDestinationFolder, err))
		}

	}
}

func TestNodeInstallCommand(t *testing.T) {

	for project_type, props := range projects {
		proj := Project{
			Name:      string(project_type),
			Directory: path.Join(root, props["directory"]),
			Resources: os.DirFS(root),
		}

		if project_type == NodeProject {
			expected := "npm i"
			installCommand := proj.NodeInstallCommand()
			if installCommand != expected {
				t.Fatalf("Expected '%s', got '%s'", expected, installCommand)
			}

			packageLockFile := path.Join(root, props["directory"], "package-lock.json")
			_, err := os.Create(packageLockFile)
			if err != nil {
				t.Errorf("Failed to create package-lock.json file for testing: %s", err)
			}

			expected = "npm ci"
			installCommand = proj.NodeInstallCommand()
			if installCommand != expected {
				t.Fatalf("Expected '%s', got '%s'", expected, installCommand)
			}
			os.Remove(packageLockFile)
		}
	}
}
