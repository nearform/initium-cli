package main

import (
	"fmt"
	"log"
	"os"
	"context"
	"time"
	"io"

	"k8s-kurated-addons.cli/src/services/docker"
    "k8s-kurated-addons.cli/src/utils/logger"

	"github.com/urfave/cli/v2"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/term"
)


func main() {
    app := &cli.App{
    		Name:  "k8s kurated addons",
    		Usage: "CLI tool ",
    		Action: func(cCtx *cli.Context) error {
    			appName := cCtx.String("app-name")
    			dockerFilePath := cCtx.String("dockerfile")
    			repoName := cCtx.String("repo-name")
    			run(appName, repoName, dockerFilePath)
    			return nil
    		},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "app-name"},
				&cli.StringFlag{Name: "repo-name"},
				&cli.StringFlag{Name: "dockerfile"},
				&cli.BoolFlag{Name: "non-interactive", Aliases: []string{"ni"}},
			},
    	}

    	if err := app.Run(os.Args); err != nil {
    		log.Fatal(err)
    	}
}

func run(appName string, repoName string, dockerFilePath string) error {
    fmt.Println("nearForm: k8s kurated addons")

    dockerService := docker.DockerService {
        DockerDirectory: dockerFilePath,
        DockerFileName: "Dockerfile",
        ContainerRepo: repoName,
        AppName: appName,
    }

    loggerUtil := logger.LoggerUtil{}
    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if (err != nil) {
        loggerUtil.PrintError("Failed to create docker client: ", err)
    }

    dockerService.Build(cli)
    dockerService.Push(cli)
    return nil
}

func buildImage(dockerClient *client.Client, dockerFilePath string, appName string, repoName string) error {

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
    defer cancel()

    buildContext, err := archive.TarWithOptions(dockerFilePath, &archive.TarOptions{})
    if (err != nil) {
        printErr("Failed to create build context", err)
    }

    buildOptions := types.ImageBuildOptions{
        Dockerfile: "Dockerfile",
        Tags:       []string{repoName + "/" + appName},
        Remove:     true,
    }

    buildResponse, err := dockerClient.ImageBuild(ctx, buildContext, buildOptions)
    if (err != nil) {
        printErr("Failed to build docker image", err)
    }

    defer buildResponse.Body.Close()

    displayMessagesStream(buildResponse.Body)

    return nil
}

func displayMessagesStream(body io.Reader) error {

    termFd, isTerm := term.GetFdInfo(os.Stdout)
    err := jsonmessage.DisplayJSONMessagesStream(body, os.Stdout, termFd, isTerm, nil)
    if (err != nil) {
        printErr("Failed to display logs", err)
    }

    return nil
}

func printErr(message string, err error) error {
    fmt.Println(message, err)
    os.Exit(1)

    return nil
}