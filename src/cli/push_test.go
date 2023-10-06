package cli

import (
	"fmt"
	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/services/project"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"gotest.tools/v3/assert"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
)

func TestShouldSendImagePushRequestToDockerDaemon(t *testing.T) {
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

	err = ds.Push()
	if err != nil {
		t.Fatalf(fmt.Sprintf("Error: %s", err))
	}

	assert.Check(t, len(apiRequests) == 2, "Expected 2 requests to be sent to Docker API")

	imageTagRequest := apiRequests[0]
	assert.Assert(t, imageTagRequest.httpMethod == http.MethodPost, "Expected POST method to be called")
	assert.Assert(t, strings.HasSuffix(imageTagRequest.url, "/images/test:v1.1.0/tag"), "Expected URL suffix to be /images/test:v1.1.0/tag")

	imagePushRequest := apiRequests[1]
	assert.Assert(t, imagePushRequest.httpMethod == http.MethodPost, "Expected POST method to be called")
	assert.Assert(t, strings.HasSuffix(imagePushRequest.url, "/images/example.org/test/push"), "Expected URL suffix to be /images/test:v1.1.0/tag")
}
