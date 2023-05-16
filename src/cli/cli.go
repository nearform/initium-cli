package cli

import (
	"embed"
	"log"
	"os"
	"path"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"

	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"

	"github.com/urfave/cli/v2"
)

type CLI struct {
	Resources embed.FS
	CWD       string
}

func (c CLI) newProject(cCtx *cli.Context) project.Project {
	return project.New(
		cCtx.String("app-name"),
		cCtx.String("project-directory"),
		cCtx.String("runtime-version"),
		cCtx.String("app-version"),
		cCtx.String("docker-image"),
		c.Resources,
	)
}

func (c CLI) Run() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "app-name",
			Usage:   "The name of the app",
			Value:   path.Base(c.CWD),
			EnvVars: []string{"KKA_APP_NAME"},
		},
		&cli.StringFlag{
			Name:    "app-version",
			Usage:   "The version of your application",
			Value:   "latest",
			EnvVars: []string{"KKA_VERSION"},
		},
		&cli.StringFlag{
			Name:    "project-directory",
			Usage:   "The directory in which your Dockerfile lives",
			Value:   defaults.ProjectDirectory,
			EnvVars: []string{"KKA_PROJECT_DIRECTORY"},
		},
		&cli.StringFlag{
			Name:    "repo-name",
			Usage:   "The base address of the container repository",
			Value:   defaults.RepoName,
			EnvVars: []string{"KKA_REPO_NAME"},
		},
		&cli.StringFlag{
			Name:    "dockerfile-name",
			Usage:   "The name of the Dockerfile",
			Value:   defaults.DockerfileName,
			EnvVars: []string{"KKA_DOCKERFILE_NAME"},
		},
	}

	app := &cli.App{
		Name:  "k8s kurated addons",
		Usage: "kka-cli",
		Flags: flags,
		Commands: []*cli.Command{
			c.BuildCMD(),
			c.PushCMD(),
			c.DeployCMD(),
			c.DeleteCMD(),
			c.OnMainCMD(),
			c.TemplateCMD(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
