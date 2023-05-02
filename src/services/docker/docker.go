package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/moby/term"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
)

type DockerService struct {
	project        project.Project
	DockerFileName string
	ContainerRepo  string
	Client         client.Client
	AuthConfig     types.AuthConfig
}

// Create a new instance of the DockerService
func New(project project.Project, dockerFileName string, containerRepo string) (DockerService, error) {
	client, err := getClient()
	if err != nil {
		return DockerService{}, err
	}

	return DockerService{
		project:        project,
		DockerFileName: dockerFileName,
		ContainerRepo:  containerRepo,
		Client:         *client,
	}, nil
}

func streamOutput(body io.ReadCloser) error {
	fd, isTerminal := term.GetFdInfo(os.Stdout)
	if err := jsonmessage.DisplayJSONMessagesStream(body, os.Stdout, fd, isTerminal, nil); err != nil {
		return fmt.Errorf("Failed to print Docker build output %v", err)
	}
	return nil
}

func getClient() (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Error("Failed to create docker client: ", err)
		return nil, err
	}

	return cli, nil
}

func (ds DockerService) RemoteTag() string {
	tag := ds.project.Version
	directory := ds.project.Directory
	if directory != defaults.ProjectDirectory {
		return fmt.Sprintf("%s/%s/%s:%s", ds.ContainerRepo, ds.project.Name, directory, tag)
	}
	return fmt.Sprintf("%s/%s:%s", ds.ContainerRepo, ds.project.Name, tag)
}

func (ds DockerService) LocalTag() string {
	tag := ds.project.Version
	directory := ds.project.Directory
	if directory != defaults.ProjectDirectory {
		return fmt.Sprintf("%s/%s:%s", ds.project.Name, directory, tag)
	}
	return fmt.Sprintf("%s:%s", ds.project.Name, tag)
}

// Build Docker image
func (ds DockerService) Build() error {
	log.SetLevel(log.DebugLevel)
	log.Infof("Building %s", ds.LocalTag())

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*5)
	defer cancel()

	// Get the context for the docker build
	buildContext, err := archive.TarWithOptions(ds.project.Directory, &archive.TarOptions{})
	if err != nil {
		return fmt.Errorf("Failed to create build context %v", err)
	}

	// Get the options for the docker build
	buildOptions := types.ImageBuildOptions{
		Dockerfile: ds.DockerFileName,
		Tags:       []string{ds.LocalTag()},
		Remove:     true,
	}

	// Build the image
	buildResponse, err := ds.Client.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		return fmt.Errorf("Failed to build docker image %v", err)
	}
	defer buildResponse.Body.Close()

	if err = streamOutput(buildResponse.Body); err != nil {
		return err
	}
	return nil
}

// Push Docker image
func (ds DockerService) Push() error {
	log.Infof("Pushing to %s", ds.RemoteTag())
	log.Debug("User: %s", ds.AuthConfig.Username)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	encodedJSON, err := json.Marshal(ds.AuthConfig)
	if err != nil {
		return err
	}
	ipo := types.ImagePushOptions{
		RegistryAuth: base64.URLEncoding.EncodeToString(encodedJSON),
	}

	err = ds.Client.ImageTag(ctx, ds.LocalTag(), ds.RemoteTag())
	if err != nil {
		return fmt.Errorf("Tagging local image for remote %v", err)
	}

	pushResponse, err := ds.Client.ImagePush(ctx, ds.RemoteTag(), ipo)
	if err != nil {
		return fmt.Errorf("Failed to push docker image %v", err)
	}
	defer pushResponse.Close()

	if err = streamOutput(pushResponse); err != nil {
		return err
	}

	return nil
}
