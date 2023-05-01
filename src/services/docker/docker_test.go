package docker

import (
	"testing"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
)

func TestLocalTag(t *testing.T) {
	ds := DockerService{
		project: project.Project{
			Directory: defaults.ProjectDirectory,
			Name:      "test",
			Version:   "v1.1.0",
		},
	}
	localTag := ds.LocalTag()
	expected := "test:v1.1.0"
	if localTag != expected {
		t.Fatalf("Expected '%s' got %s", expected, localTag)
	}

	ds.project.Directory = "Subproject"

	localTag = ds.LocalTag()
	expected = "test/Subproject:v1.1.0"
	if localTag != expected {
		t.Fatalf("Expected '%s' got %s", expected, localTag)
	}
}

func TestRemoteTag(t *testing.T) {
	ds := DockerService{
		ContainerRepo: "example.org",
		project: project.Project{
			Directory: defaults.ProjectDirectory,
			Name:      "test",
			Version:   "v1.1.0",
		},
	}

	localTag := ds.RemoteTag()
	expected := "example.org/test:v1.1.0"
	if localTag != expected {
		t.Fatalf("Expected '%s' got %s", expected, localTag)
	}

	ds.project.Directory = "Subproject"

	localTag = ds.RemoteTag()
	expected = "example.org/test/Subproject:v1.1.0"
	if localTag != expected {
		t.Fatalf("Expected '%s' got %s", expected, localTag)
	}
}
