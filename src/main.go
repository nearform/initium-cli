package main

import (
	"log"
	"os"

	"k8s-kurated-addons.cli/src/services/docker"

    "k8s-kurated-addons.cli/src/utils/logger"
    "k8s-kurated-addons.cli/src/utils/helper"
    "k8s-kurated-addons.cli/src/utils/defaults"

	"github.com/urfave/cli/v2"
)

func main() {
    app := &cli.App{
    		Name:  "k8s kurated addons",
    		Usage: "CLI tool for KKA",
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
				&cli.StringFlag{Name: "dockerfile-name"},
			},
    	}

    	if err := app.Run(os.Args); err != nil {
    		log.Fatal(err)
    	}
}

// Initialize default values
func initArguments(appName *string, repoName *string, dockerFilePath *string, dockerFileName *string) {
    helper.ReplaceIfEmpty(appName, defaults.DefaultAppName)
    helper.ReplaceIfEmpty(dockerFilePath, defaults.DefaultDockerDirectory)
    helper.ReplaceIfEmpty(repoName, defaults.DefaultRepoName)
    helper.ReplaceIfEmpty(dockerFileName, defaults.DefaultDockerfileName)
}

// Run the CLI
func run(appName string, repoName string, dockerFilePath string, dockerFileName string) error {
    logger.PrintInfo("nearForm: k8s kurated addons CLI")
    logger.PrintInfo("Dockerfile Location: " + dockerFilePath + "/" + dockerFileName)
    logger.PrintInfo("Building to: " + repoName + "/" + appName)

    dockerService := docker.New(dockerFilePath, dockerFileName, repoName, appName)

    dockerService.Build()
    dockerService.Push()

    return nil
}
