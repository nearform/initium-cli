package docker

import (
    	"context"
    	"time"
    	"k8s-kurated-addons.cli/src/utils/logger"
    	"github.com/docker/docker/api/types"
    	"github.com/docker/docker/client"
    	"github.com/docker/docker/pkg/archive"
)

type DockerService struct {
    DockerDirectory string
    DockerFileName string
    ContainerRepo string
    AppName string
}

func (ds DockerService) Build(client *client.Client) error {

        loggerUtil := logger.LoggerUtil{}
        loggerUtil.PrintInfo("Building...")
        ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
        defer cancel()

        buildContext, err := archive.TarWithOptions(ds.DockerDirectory, &archive.TarOptions{})
        if (err != nil) {
            loggerUtil.PrintError("Failed to create build context", err)
        }

        buildOptions := types.ImageBuildOptions{
            Dockerfile: ds.DockerFileName,
            Tags:       []string{ds.ContainerRepo + "/" + ds.AppName},
            Remove:     true,
        }

        loggerUtil.PrintInfo(ds.DockerDirectory + ds.DockerFileName + ds.ContainerRepo + ds.AppName)


        buildResponse, err := client.ImageBuild(ctx, buildContext, buildOptions)
        loggerUtil.PrintInfo(ds.ContainerRepo + "/" + ds.AppName)
        if (err != nil) {
            loggerUtil.PrintError("Failed to build docker image", err)
        }

        defer buildResponse.Body.Close()

        loggerUtil.PrintStream(buildResponse.Body)

        return nil
}

func (ds DockerService) Push(client *client.Client) error {
    loggerUtil := logger.LoggerUtil{}
    loggerUtil.PrintInfo("Pushing...")

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
    defer cancel()

    pushResponse, err := client.ImagePush(ctx, ds.ContainerRepo + "/" + ds.AppName, types.ImagePushOptions{})
    if (err != nil) {
        loggerUtil.PrintError("Failed to push docker image", err)
    }

    defer pushResponse.Close()
    loggerUtil.PrintStream(pushResponse)

    return nil
}
