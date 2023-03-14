package main

import (
	"fmt"
	"log"
	"os"

	"k8s-kurated-addons.cli/src/services/docker"
    "k8s-kurated-addons.cli/src/utils/logger"

	"github.com/urfave/cli/v2"
	"github.com/docker/docker/client"
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