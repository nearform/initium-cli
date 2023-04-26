package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"k8s-kurated-addons.cli/src/services/project"
	"k8s-kurated-addons.cli/src/utils/logger"
)

type DockerService struct {
	project        project.Project
	DockerFileName string
	ContainerRepo  string
	Client         client.Client
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

// Build Docker image
func (ds DockerService) Build() error {
	logger.PrintInfo("Building...")

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
		Tags:       []string{ds.ContainerRepo + "/" + ds.project.Name},
		Remove:     true,
	}

	// Build the image
	buildResponse, err := ds.Client.ImageBuild(ctx, buildContext, buildOptions)
	logger.PrintInfo(ds.ContainerRepo + "/" + ds.project.Name)
	if err != nil {
		logger.PrintError("Failed to build docker image", err)
	}

	defer buildResponse.Body.Close()

	logger.PrintStream(buildResponse.Body)

	return nil
}

// Push Docker image
func (ds DockerService) Push() error {
	logger.PrintInfo("Pushing...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
	defer cancel()

	// Push the image
	pushResponse, err := ds.Client.ImagePush(ctx, ds.ContainerRepo+"/"+ds.project.Name, types.ImagePushOptions{})
	if err != nil {
		logger.PrintError("Failed to push docker image", err)
	}

	defer pushResponse.Close()
	logger.PrintStream(pushResponse)

	return nil
}
