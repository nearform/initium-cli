package project

import (
	"fmt"
	"os"
	"path"
	"testing"
	"strings"
)

var projects = map[string]map[string]string{
	"node":    {"directory": "example"},
	"go":      {"directory": "."},
	"invalid": {"directory": "src"},
}

var root = "../../../"

func TestDetectType(t *testing.T) {

	for project_type, props := range projects {
		test_proj_type := Project{Name: project_type,
			Directory: path.Join(root, props["directory"])}

		proj_type, err := test_proj_type.detectType()

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
		proj_dockerfile := Project{Name: project_type,
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
	
	proj_runtime := Project{Name: "test",
		Directory: path.Join(root, projects["node"]["directory"]),
		Resources: os.DirFS(root),
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