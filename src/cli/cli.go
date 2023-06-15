package cli

import (
	"embed"
	"os"
	"sort"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"

	"github.com/charmbracelet/log"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"
	"github.com/urfave/cli/v2"
)

type CLI struct {
	Resources     embed.FS
	CWD           string
	DockerService docker.DockerService
	Logger        *log.Logger
	project       project.Project
	dockerImage   docker.DockerImage
}

func (c *CLI) init(cCtx *cli.Context) {
	repoName := cCtx.String("repo-name")
	dockerFileName := cCtx.String("dockerfile-name")
	appName := cCtx.String("app-name")
	version := cCtx.String("app-version")
	projectDirectory := cCtx.String("project-directory")

	project := project.New(
		appName,
		projectDirectory,
		cCtx.String("runtime-version"),
		version,
		c.Resources,
	)

	dockerImage := docker.DockerImage{
		Registry:  repoName,
		Name:      appName,
		Directory: projectDirectory,
		Tag:       version,
	}

	dockerService, err := docker.New(project, dockerImage, dockerFileName)
	if err != nil {
		logger.PrintError("Error creating docker service", err)
	}

	c.DockerService = dockerService
	c.dockerImage = dockerImage
	c.project = project
}

func (c *CLI) getProject(cCtx *cli.Context) *project.Project {
	if (c.project == project.Project{}) {
		c.init(cCtx)
	}
	return &c.project
}

func (c CLI) Run() {
	app := &cli.App{
		Name:  "k8s kurated addons",
		Usage: "kka-cli",
		Flags: c.CommandFlags(App),
		Commands: []*cli.Command{
			c.BuildCMD(),
			c.PushCMD(),
			c.DeployCMD(),
			c.DeleteCMD(),
			c.OnMainCMD(),
			c.TemplateCMD(),
			c.InitCMD(),
		},
		Before: func(ctx *cli.Context) error {
			err := c.loadFlagsFromConfig(ctx)

			if err != nil {
				c.Logger.Debug("failed to load config", err)
			}

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
