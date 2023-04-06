package kkacli

import (
    "log"
    "os"

    "k8s-kurated-addons.cli/src/services/docker"
    "k8s-kurated-addons.cli/src/services/kubernetes"
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
            appPort := cCtx.String("app-port")
            dockerFilePath := cCtx.String("dockerfile-directory")
            repoName := cCtx.String("repo-name")
            dockerFileName := cCtx.String("dockerfile-name")
            manifestPath := cCtx.String("manifest-path")

            localService := local.LocalService{
                HasDockerfile: dockerFileName != "",
                HasManifest: manifestPath != "",
            }

            parameters := initParameters(&appName, &appPort, &repoName, &dockerFilePath, &dockerFileName, &manifestPath, localService)
            runCli(parameters)
//             cleanUp(localService, dockerFilePath)

            return nil
        },
        Flags: []cli.Flag{
            &cli.StringFlag{Name: "app-name", Usage: "The name of the app"},
            &cli.StringFlag{Name: "app-port", Usage: "The name of the app"},
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
func initParameters(appName *string, appPort *string, repoName *string, dockerFilePath *string, dockerFileName *string, manifestPath *string, localService local.LocalService) helper.Parameters {

    helper.ReplaceIfEmpty(dockerFilePath, defaults.DefaultDockerDirectory)

    // If no Dockerfile is given, we will create one
    if (!localService.HasDockerfile) {
        localService.CreateDockerfile(*dockerFilePath)
        helper.ReplaceIfEmpty(dockerFileName, defaults.DefaultDockerfileName)
    }

    helper.ReplaceIfEmpty(appName, defaults.DefaultAppName)
    helper.ReplaceIfEmpty(appPort, defaults.DefaultAppPort)
    helper.ReplaceIfEmpty(repoName, defaults.DefaultRepoName)

    // If no manifest is given, we will create one
    if (!localService.HasManifest) {
        localService.CreateManifest(*appPort, *appName, *repoName)
        helper.ReplaceIfEmpty(manifestPath, defaults.DefaultManifestPath)
    }

    return helper.Parameters{
        AppName: *appName,
        AppPort: *appPort,
        RepoName: *repoName,
        DockerFilePath: *dockerFilePath,
        DockerFileName: *dockerFileName,
        ManifestPath: *manifestPath,
    }
}

// Run the CLI
func runCli(parameters helper.Parameters) error {
    logger.PrintInfo("Dockerfile Location: " + parameters.DockerFilePath + "/" + parameters.DockerFileName)
    logger.PrintInfo("Building to: " + parameters.RepoName + "/" + parameters.AppName)

    dockerService := docker.New(parameters.DockerFilePath, parameters.DockerFileName, parameters.RepoName, parameters.AppName)

    dockerService.Build()
    dockerService.Push()

    kubernetesService := kubernetes.New()
    kubernetesService.DeployManifest(parameters.ManifestPath)

    return nil
}

// Clean up any temporary files
func cleanUp(localService local.LocalService, dockerFilePath string) {
    if (!localService.HasDockerfile) {
        localService.RemoveDockerfile(dockerFilePath)
    }
}
