package cli

import (
	"io"
	"io/fs"

	"os"
	"path/filepath"
	"sort"

	"github.com/nearform/initium-cli/src/services/project"
	"k8s.io/utils/strings/slices"

	"github.com/charmbracelet/log"
	"github.com/nearform/initium-cli/src/services/docker"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"github.com/nearform/initium-cli/src/utils/logger"
	"github.com/urfave/cli/v2"
)

type icli struct {
	Resources     fs.FS
	CWD           string
	DockerService docker.DockerService
	Logger        *log.Logger
	project       project.Project
	dockerImage   docker.DockerImage
	flags         flags
	Writer        io.Writer
}

func NewWithOptions(resources fs.FS, logger *log.Logger, writer io.Writer) icli {
	return icli{
		Resources: resources,
		Logger:    logger,
		Writer:    writer,
		flags:     InitFlags(),
	}
}

func New(resources fs.FS) icli {
	return NewWithOptions(
		resources,
		log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.ParseLevel(os.Getenv("INITIUM_LOG_LEVEL")),
			ReportCaller:    true,
			ReportTimestamp: true,
		}),
		os.Stdout,
	)
}

func (c icli) baseBeforeFunc(ctx *cli.Context) error {
	if err := c.loadFlagsFromConfig(ctx); err != nil {
		return err
	}

	if err := c.checkRequiredFlags(ctx, []string{}); err != nil {
		return err
	}
	return nil
}

func (c *icli) init(cCtx *cli.Context) error {
	appName := cCtx.String(appNameFlag)
	version := cCtx.String(appVersionFlag)
	projectLanguage := cCtx.String(projectLanguageFlag)
	projectDirectory := cCtx.String(projectDirectoryFlag)
	absProjectDirectory, err := filepath.Abs(cCtx.String(projectDirectoryFlag))
	registry := cCtx.String(repoNameFlag)
	dockerFileName := cCtx.String(dockerFileNameFlag)

	if dockerFileName == "" {
		dockerFileName = defaults.GeneratedDockerFile
	}

	if err != nil {
		c.Logger.Warnf("could not get abs of %s", projectDirectory)
		absProjectDirectory = projectDirectory
	}

	project := project.New(
		appName,
		projectLanguage,
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
		Registry:  registry,
		Name:      dockerImageName,
		Directory: absProjectDirectory,
		Tag:       version,
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

func (c *icli) getProject(cCtx *cli.Context) (*project.Project, error) {
	if (c.project == project.Project{}) {
		err := c.init(cCtx)
		if err != nil {
			return nil, err
		}
	}
	return &c.project, nil
}

func (c icli) Run(args []string) error {
	app := &cli.App{
		Name:  "initium",
		Usage: "icli of the Initium project",
		Flags: c.CommandFlags([]FlagsType{App}),
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
