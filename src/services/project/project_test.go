package project

import (
	"testing"
	"fmt"

	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
)

func TestDetectType(t *testing.T) {

	root := defaults.RootDirectoryTests
	projs := []map[string]string{
		{"name": "node", "directory": "example"},
		{"name": "go", "directory": "."},
	}

	for _, value := range projs {
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