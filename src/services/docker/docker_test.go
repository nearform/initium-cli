package docker

import (
	"testing"

	"github.com/nearform/initium-cli/src/services/project"
	"github.com/nearform/initium-cli/src/utils/defaults"
)

func TestLocalTag(t *testing.T) {

	dockerImage := DockerImage{
		Directory: defaults.ProjectDirectory,
		Name:      "test",
		Tag:       "v1.1.0",
	}
	ds := DockerService{
		project: project.Project{
			Directory: defaults.ProjectDirectory,
			Name:      "test",
			Version:   "v1.1.0",
		},
		dockerImage: dockerImage,
	}
	localTag := ds.dockerImage.LocalTag()
	expected := "test:v1.1.0"
	if localTag != expected {
		t.Fatalf("Expected '%s' got %s", expected, localTag)
	}
}

func TestRemoteTag(t *testing.T) {
	dockerImage := DockerImage{
		Registry:  "example.org",
		Directory: defaults.ProjectDirectory,
		Name:      "test",
		Tag:       "v1.1.0",
	}

	ds := DockerService{
		project: project.Project{
			Directory: defaults.ProjectDirectory,
			Name:      "test",
			Version:   "v1.1.0",
		},
		dockerImage: dockerImage,
	}

	remoteTag := ds.dockerImage.RemoteTag()
	expected := "example.org/test:v1.1.0"
	if remoteTag != expected {
		t.Fatalf("Expected '%s' got %s", expected, remoteTag)
	}
}
