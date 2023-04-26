package kkacli

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"

	"k8s-kurated-addons.cli/src/services/docker"

	"k8s-kurated-addons.cli/src/utils/defaults"
	"k8s-kurated-addons.cli/src/utils/logger"

	"github.com/urfave/cli/v2"
)

//go:embed dockerfiles
var res embed.FS

func Run() {
	logger.PrintInfo("nearForm: k8s kurated addons CLI")
	app := &cli.App{
		Name:  "k8s kurated addons",
		Usage: "kka-cli",
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "build a container image from the project directory",
				Action: func(cCtx *cli.Context) error {
					appName := cCtx.String("app-name")
					projectPath := cCtx.String("project-directory")
					repoName := cCtx.String("repo-name")
					dockerFileName := cCtx.String("dockerfile-name")

					runCli(appName, repoName, projectPath, dockerFileName, "build")

					return nil
				},
			},
			{
				Name:  "push",
				Usage: "push the container image ot a registry",
				Action: func(cCtx *cli.Context) error {
					appName := cCtx.String("app-name")
					projectPath := cCtx.String("project-directory")
					repoName := cCtx.String("repo-name")
					dockerFileName := cCtx.String("dockerfile-name")

					runCli(appName, repoName, projectPath, dockerFileName, "push")
					return nil
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "app-name",
				Usage:   "The name of the app",
				Value:   defaults.AppName,
				EnvVars: []string{"KKA_APP_NAME"},
			},
			&cli.StringFlag{
				Name:    "repo-name",
				Usage:   "The base address of the container repository you are wanting to push the image to.",
				Value:   defaults.RepoName,
				EnvVars: []string{"KKA_REPO_NAME"},
			},
			&cli.StringFlag{
				Name:    "project-directory",
				Usage:   "The directory in which your Dockerfile lives.",
				Value:   defaults.ProjectDirectory,
				EnvVars: []string{"KKA_PROJECT_DIRECTORY"},
			},
			&cli.StringFlag{
				Name:    "dockerfile-name",
				Usage:   "The name of the Dockerfile",
				Value:   defaults.DockerfileName,
				EnvVars: []string{"KKA_DOCKERFILE_NAME"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func detectProjectType(projectDirectory string) (string, error) {
	if _, err := os.Stat(path.Join(projectDirectory, "package.json")); err == nil {
		return "node", err
	} else {
		return "", fmt.Errorf("Cannot detect project type")
	}
}

// Run the CLI
func runCli(appName string, repoName string, projectDirectory string, dockerFileName string, action string) error {
	logger.PrintInfo("Dockerfile Location: " + path.Join(projectDirectory, dockerFileName))
	logger.PrintInfo("Pushing to: " + repoName + "/" + appName)
	dockerService := docker.New(projectDirectory, dockerFileName, repoName, appName)

	if action == "push" {
		dockerService.Push()
		return nil
	}

	defer docker.DeleteDockerFile(projectDirectory)
	projectType, err := detectProjectType(projectDirectory)
	if err != nil {
		return fmt.Errorf("Detect project type: %v", err)
	}

	dockerfileTemplate := path.Join(".", "dockerfiles", "Dockerfile."+projectType)
	dockerfileContent, err := res.ReadFile(dockerfileTemplate)
	if err != nil {
		return fmt.Errorf("Getting dockerfile content: %v", err)
	}

	err = docker.PersistDockerFile(projectDirectory, dockerfileContent)
	if err != nil {
		return fmt.Errorf("Persisting docker file content: %v", err)
	}

	dockerService.Build()
	return nil
}
