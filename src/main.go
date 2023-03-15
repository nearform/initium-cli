package main

import (
	"log"
	"os"

	"k8s-kurated-addons.cli/src/services/docker"
    "k8s-kurated-addons.cli/src/utils/logger"
    "k8s-kurated-addons.cli/src/utils/helper"

	"github.com/urfave/cli/v2"
	"github.com/docker/docker/client"
)


func main() {
    app := &cli.App{
    		Name:  "k8s kurated addons",
    		Usage: "CLI tool ",
    		Action: func(cCtx *cli.Context) error {
                appName := cCtx.String("app-name")
       			dockerFilePath := cCtx.String("dockerfile-directory")
       			repoName := cCtx.String("repo-name")
       			dockerFileName := cCtx.String("dockerfile-name")
       			initArguments(&appName, &repoName, &dockerFilePath, &dockerFileName)
    			run(appName, repoName, dockerFilePath, dockerFileName)
    			return nil
    		},
			Flags: []cli.Flag{
				&cli.StringFlag{Name: "app-name"},
				&cli.StringFlag{Name: "repo-name"},
				&cli.StringFlag{Name: "dockerfile-directory"},
				&cli.BoolFlag{Name: "non-interactive", Aliases: []string{"ni"}},
			},
    	}

    	if err := app.Run(os.Args); err != nil {
    		log.Fatal(err)
    	}
}

func initArguments(appName *string, repoName *string, dockerFilePath *string, dockerFileName *string) {
    helper.ReplaceIfEmpty(appName, "k8s-kurated-addons-cli")
    helper.ReplaceIfEmpty(dockerFilePath, "docker")
    helper.ReplaceIfEmpty(repoName, "ghcr.io/nearform")
    helper.ReplaceIfEmpty(dockerFileName, "Dockerfile")
}

func run(appName string, repoName string, dockerFilePath string, dockerFileName string) error {
    loggerUtil := logger.LoggerUtil{}
    loggerUtil.PrintInfo("nearForm: k8s kurated addons")
    loggerUtil.PrintInfo("Dockerfile Location: " + dockerFilePath + "/" + dockerFileName)
    loggerUtil.PrintInfo("Building to: " + repoName + "/" + appName)

    dockerService := docker.DockerService {
        DockerDirectory: dockerFilePath,
        DockerFileName: "Dockerfile",
        ContainerRepo: repoName,
        AppName: appName,
    }

    cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
    if (err != nil) {
        loggerUtil.PrintError("Failed to create docker client: ", err)
    }

    dockerService.Build(cli)
    dockerService.Push(cli)

    return nil
}