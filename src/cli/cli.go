package cli

import (
	"embed"
	"io"
	"os"
	"path/filepath"
	"sort"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"
	"k8s.io/utils/strings/slices"

	"github.com/charmbracelet/log"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
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
	Writer        io.Writer
}

func (c CLI) baseBeforeFunc(ctx *cli.Context) error {
	if err := c.loadFlagsFromConfig(ctx); err != nil {
		return err
	}

	if err := c.checkRequiredFlags(ctx, []string{}); err != nil {
		return err
	}
	return nil
}

func (c *CLI) init(cCtx *cli.Context) error {
	appName := cCtx.String(appNameFlag)
	version := cCtx.String(appVersionFlag)
	projectDirectory := cCtx.String(projectDirectoryFlag)
	absProjectDirectory, err := filepath.Abs(cCtx.String(projectDirectoryFlag))

	if err != nil {
		c.Logger.Warnf("could not get abs of %s", projectDirectory)
		absProjectDirectory = projectDirectory
	}

	project := project.New(
		appName,
		projectDirectory,
		cCtx.String(runtimeVersionFlag),
		version,
		c.Resources,
	)

	dockerImageName := appName
	invalidBases := []string{".", string(os.PathSeparator)}
	if projectDirectory != defaults.ProjectDirectory {
		base := filepath.Base(absProjectDirectory)
		if !slices.Contains(invalidBases, base) && base != dockerImageName {
			dockerImageName = appName + "/" + base
		}
	}

	dockerImage := docker.DockerImage{
		Registry:  cCtx.String(repoNameFlag),
		Name:      dockerImageName,
		Directory: absProjectDirectory,
		Tag:       version,
	}

	dockerFileName := cCtx.String(dockerFileNameFlag)
	if dockerFileName == "" {
		dockerFileName = defaults.GeneratedDockerFile
	}

	dockerService, err := docker.New(project, dockerImage, dockerFileName)
	if err != nil {
		logger.PrintError("Error creating docker service", err)
	}

	c.DockerService = dockerService
	c.dockerImage = dockerImage
	c.project = project
	return nil
}

func (c *CLI) getProject(cCtx *cli.Context) (*project.Project, error) {
	if (c.project == project.Project{}) {
		err := c.init(cCtx)
		if err != nil {
			return nil, err
		}
	}
	return &c.project, nil
}

func (c CLI) Run(args []string) error {
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
			c.OnBranchCMD(),
			c.TemplateCMD(),
			c.InitCMD(),
		},
		Before: func(ctx *cli.Context) error {
			if err := c.loadFlagsFromConfig(ctx); err != nil {
				return err
			}

			if err := c.checkRequiredFlags(ctx, []string{}); err != nil {
				return err
			}
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app.Run(args)
}
