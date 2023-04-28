package kkacli

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/docker/docker/api/types"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	knative "github.com/nearform/k8s-kurated-addons-cli/src/services/k8s"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"

	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"

	"github.com/urfave/cli/v2"
)

//go:embed assets
var res embed.FS

func newProject(cCtx *cli.Context) project.Project {
	return project.New(
		cCtx.String("app-name"),
		cCtx.String("project-directory"),
		cCtx.String("runtime-version"),
		cCtx.String("version"),
		res,
	)
}

// Run the CLI
func build(docker docker.DockerService, project project.Project) error {
	logger.PrintInfo("Dockerfile Location: " + path.Join(project.Directory, docker.DockerFileName))

	defer project.DeleteDockerFile()
	err := project.AddDockerFile()
	if err != nil {
		return fmt.Errorf("Persisting docker file content: %v", err)
	}

	return docker.Build()
}

func Run() {
	logger.PrintInfo("nearForm: k8s kurated addons CLI")
	app := &cli.App{
		Name:  "k8s kurated addons",
		Usage: "kka-cli",
		Commands: []*cli.Command{
			{
				Name:  "build",
				Usage: "build a container image from the project directory",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "runtime-version",
						EnvVars: []string{"KKA_RUNTIME_VERSION"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					repoName := cCtx.String("repo-name")
					dockerFileName := cCtx.String("dockerfile-name")
					project := newProject(cCtx)
					docker := docker.New(project, dockerFileName, repoName)
					return build(docker, project)
				},
			},
			{
				Name:  "push",
				Usage: "push the container image to a registry",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "registry-user",
						EnvVars: []string{"KKA_REGISTRY_USER"},
					},
					&cli.StringFlag{
						Name:    "registry-password",
						EnvVars: []string{"KKA_REGISTRY_PASSWORD"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					repoName := cCtx.String("repo-name")
					dockerFileName := cCtx.String("dockerfile-name")
					project := newProject(cCtx)
					docker := docker.New(project, dockerFileName, repoName)
					docker.AuthConfig = types.AuthConfig{
						Username: cCtx.String("registry-user"),
						Password: cCtx.String("registry-password"),
					}
					return docker.Push()
				},
			},
			{
				Name:  "deploy",
				Usage: "deploy the application as a knative service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "endpoint",
						EnvVars: []string{"KKA_ENDPOINT"},
					},
					&cli.StringFlag{
						Name:    "token",
						EnvVars: []string{"KKA_TOKEN"},
					},
					&cli.StringFlag{
						Name:    "ca-crt",
						EnvVars: []string{"KKA_CA_CERT"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					config, err := knative.Config(
						cCtx.String("endpoint"),
						cCtx.String("token"),
						[]byte(cCtx.String("ca-crt")),
					)
					if err != nil {
						return err
					}
					project := newProject(cCtx)
					return knative.Apply(config, project)
				},
			},
			{
				Name:  "delete",
				Usage: "delete the knative service",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "endpoint",
						EnvVars: []string{"KKA_ENDPOINT"},
					},
					&cli.StringFlag{
						Name:    "token",
						EnvVars: []string{"KKA_TOKEN"},
					},
					&cli.StringFlag{
						Name:    "ca-crt",
						EnvVars: []string{"KKA_CA_CERT"},
					},
				},
				Action: func(cCtx *cli.Context) error {
					config, err := knative.Config(
						cCtx.String("endpoint"),
						cCtx.String("token"),
						[]byte(cCtx.String("ca-crt")),
					)
					if err != nil {
						return err
					}
					project := newProject(cCtx)
					return knative.Clean(config, project)
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
			&cli.StringFlag{
				Name:    "version",
				Usage:   "The version of your application",
				Value:   "",
				EnvVars: []string{"KKA_VERSION"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
