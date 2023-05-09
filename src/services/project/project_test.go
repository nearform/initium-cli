package project

import (
	"testing"
	"fmt"
	"os"
	// _"embed"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
)

var projects = []map[string]string{
	{"name": "node", "directory": "example"},
	{"name": "go", "directory": "."},
}

var root = defaults.RootDirectoryTests

func TestDetectType(t *testing.T) {

	for _, value := range projects {
		test_proj_type := Project{Name: value["name"], 
							      Directory: fmt.Sprintf("%s%s", root, value["directory"])}

		proj_type, err := test_proj_type.detectType()

		if err != nil {
			t.Fatalf(fmt.Sprintf("Error: %s", err))
		}

		if proj_type != value["name"] {
			t.Fatalf("Error: %s project not found", value["name"])
		}
	}	
}

func TestTemplateFile(t *testing.T){

	for _, value := range projects {
		if _, err := os.Stat(fmt.Sprintf("%s/assets/docker/Dockerfile.%s.tmpl", root, value["name"])); err != nil {
			t.Fatalf("Error: Dockerfile.%s.tmpl template not found", value["name"])
		}
	}
}

func TestLoadDockerfile(t *testing.T){
	for _, value := range projects {
		proj_dockerfile := Project{Name: value["name"], 
							      Directory: fmt.Sprintf("%s%s", root, value["directory"]),
								  Resources: defaults.TemplateFS}
		_, err := proj_dockerfile.loadDockerfile()

		if err != nil {
			fmt.Println(err)
			t.Fatalf(fmt.Sprintf("Error: %s", err))
		}
		
	}
}