package docker

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"
)

type DockerService struct {
	project        project.Project
	DockerFileName string
	ContainerRepo  string
	Client         client.Client
	AuthConfig     types.AuthConfig
}

// Create a new instance of the DockerService
func New(project project.Project, dockerFileName string, containerRepo string) DockerService {
	return DockerService{
		project:        project,
		DockerFileName: dockerFileName,
		ContainerRepo:  containerRepo,
		Client:         *getClient(),
	}
}

func getClient() *client.Client {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		logger.PrintError("Failed to create docker client: ", err)
	}

	return cli
}

func (ds DockerService) Remote() string {
	tag := "latest"
	if ds.project.Version != "" {
		tag = ds.project.Version
	}
	return fmt.Sprintf("%s/%s:%s", ds.ContainerRepo, ds.project.Name, tag)
}

// Build Docker image
func (ds DockerService) Build() error {
	logger.PrintInfo("Building " + ds.Remote())

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	// Get the context for the docker build
	buildContext, err := archive.TarWithOptions(ds.project.Directory, &archive.TarOptions{})
	if err != nil {
		logger.PrintError("Failed to create build context", err)
	}

	// Get the options for the docker build
	buildOptions := types.ImageBuildOptions{
		Dockerfile: ds.DockerFileName,
		Tags:       []string{ds.Remote()},
		Remove:     true,
	}

	// Build the image
	buildResponse, err := ds.Client.ImageBuild(ctx, buildContext, buildOptions)
	if err != nil {
		logger.PrintError("Failed to build docker image", err)
	}

	defer buildResponse.Body.Close()

	logger.PrintStream(buildResponse.Body)

	return nil
}

// Push Docker image
func (ds DockerService) Push() error {
	logger.PrintInfo("Pushing to " + ds.Remote())
	logger.PrintInfo("User: " + ds.AuthConfig.Username)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	encodedJSON, err := json.Marshal(ds.AuthConfig)
	if err != nil {
		return err
	}
	ipo := types.ImagePushOptions{
		RegistryAuth: base64.URLEncoding.EncodeToString(encodedJSON),
	}

	pushResponse, err := ds.Client.ImagePush(ctx, ds.Remote(), ipo)
	defer pushResponse.Close()
	if err != nil {
		logger.PrintError("Failed to push docker image", err)
	}

	return logger.PrintStream(pushResponse)
}
