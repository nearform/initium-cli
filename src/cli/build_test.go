package cli

import (
	"fmt"
	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/services/project"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"os"
	"path"
	"testing"
)

func TestShouldSendBuildRequestToDockerDaemon(t *testing.T) {
	proj := project.Project{
		Name:      "initium-cli",
		Directory: path.Join("../../", "."),
		Resources: os.DirFS("../../"),
	}

	dockerImage := docker.DockerImage{
		Registry:  "example.org",
		Directory: defaults.ProjectDirectory,
		Name:      "test",
		Tag:       "v1.1.0",
	}

	var apiRequests []dockerApiRequest
	dockerClient, err := newAlwaysOkMockDockerClient(&apiRequests)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	ds, err := docker.NewWithDockerClient(proj, dockerImage, defaults.GeneratedDockerFile, &dockerClient)
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	err = ds.Build()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}
}
