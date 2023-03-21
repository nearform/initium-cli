package kkacli

import (
    "log"
    "os"

    "k8s-kurated-addons.cli/src/services/docker"
    "k8s-kurated-addons.cli/src/services/local"

    "k8s-kurated-addons.cli/src/utils/logger"
    "k8s-kurated-addons.cli/src/utils/helper"
    "k8s-kurated-addons.cli/src/utils/defaults"

    "github.com/urfave/cli/v2"
)

func Run() {
    logger.PrintInfo("nearForm: k8s kurated addons CLI")
    app := &cli.App{
        Name:  "k8s kurated addons",
        Usage: "kka-cli",
        Action: func(cCtx *cli.Context) error {
            appName := cCtx.String("app-name")
            dockerFilePath := cCtx.String("dockerfile-directory")
            repoName := cCtx.String("repo-name")
            dockerFileName := cCtx.String("dockerfile-name")

            localService := local.LocalService{
                HasDockerfile: dockerFileName != "",
            }

            initArguments(&appName, &repoName, &dockerFilePath, &dockerFileName, localService)
            runCli(appName, repoName, dockerFilePath, dockerFileName)
            cleanUp(localService, dockerFilePath)

            return nil
        },
        Flags: []cli.Flag{
            &cli.StringFlag{Name: "app-name", Usage: "The name of the app"},
            &cli.StringFlag{Name: "repo-name", Usage: "The base address of the container repository you are wanting to push the image to."},
            &cli.StringFlag{Name: "dockerfile-directory", Usage: "The directory in which your Dockerfile lives."},
            &cli.StringFlag{Name: "dockerfile-name", Usage: "The name of the Dockerfile"},
        },
    }

    if err := app.Run(os.Args); err != nil {
        log.Fatal(err)
    }
}

// Initialize default values
func initArguments(appName *string, repoName *string, dockerFilePath *string, dockerFileName *string, localService local.LocalService) {

    helper.ReplaceIfEmpty(dockerFilePath, defaults.DefaultDockerDirectory)

    // If no Dockerfile is given, we will create one
    if (!localService.HasDockerfile) {
        localService.CreateDockerfile(*dockerFilePath)
        helper.ReplaceIfEmpty(dockerFileName, defaults.DefaultDockerfileName)
    }

    helper.ReplaceIfEmpty(appName, defaults.DefaultAppName)
    helper.ReplaceIfEmpty(repoName, defaults.DefaultRepoName)
}

// Run the CLI
func runCli(appName string, repoName string, dockerFilePath string, dockerFileName string) error {
    logger.PrintInfo("Dockerfile Location: " + dockerFilePath + "/" + dockerFileName)
    logger.PrintInfo("Building to: " + repoName + "/" + appName)

    dockerService := docker.New(dockerFilePath, dockerFileName, repoName, appName)

    dockerService.Build()
    dockerService.Push()

    return nil
}

// Clean up any temporary files
func cleanUp(localService local.LocalService, dockerFilePath string) {
    if (!localService.HasDockerfile) {
        localService.RemoveDockerfile(dockerFilePath)
    }
}
