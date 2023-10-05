package cli

import (
	"bytes"
	"fmt"
	"github.com/docker/docker/client"
	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/services/project"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"testing"
)

type transportFunc func(*http.Request) (*http.Response, error)

func (tf transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return tf(req)
}

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

	dockerClient, err := getMockDockerClient(t)
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

func getMockDockerClient(t *testing.T) (client.Client, error) {
	handler := func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodPost {
			t.Fatalf("POST method call expected")
		}

		if !strings.HasSuffix(request.URL.Path, "/build") {
			t.Fatalf("/build URL call expected")
		}

		header := http.Header{}
		header.Set("Content-Type", "application/json")

		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("{}"))),
			Header:     header,
		}, nil
	}

	mockClient, err := client.NewClientWithOpts(client.WithHTTPClient(&http.Client{
		Transport: transportFunc(handler),
	}))
	if err != nil {
		return client.Client{}, err
	}

	return *mockClient, nil
}
