package cli

import (
	"embed"
	"log"
	"os"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"

	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"

	"github.com/urfave/cli/v2"
)

type CLI struct {
	Resources embed.FS
}

func (c CLI) newProject(cCtx *cli.Context) project.Project {
	return project.New(
		cCtx.String("app-name"),
		cCtx.String("project-directory"),
		cCtx.String("runtime-version"),
		cCtx.String("version"),
		c.Resources,
	)
}

func (c CLI) Run() {
	logger.PrintInfo("nearForm: k8s kurated addons CLI")
	app := &cli.App{
		Name:  "k8s kurated addons",
		Usage: "kka-cli",
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
		Commands: []*cli.Command{
			c.Build(),
			c.Push(),
			c.Deploy(),
			c.Delete(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
