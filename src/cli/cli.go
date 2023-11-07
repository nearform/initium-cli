package cli

import (
	"fmt"
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

type Release struct {
	Version string
	Commit  string
	Date    string
}

type icli struct {
	Resources     fs.FS
	CWD           string
	DockerService docker.DockerService
	Logger        *log.Logger
	project       *project.Project
	dockerImage   docker.DockerImage
	flags         flags
	Writer        io.Writer
	release       Release
}

func NewWithOptions(resources fs.FS, logger *log.Logger, writer io.Writer, release Release) icli {
	return icli{
		Resources: resources,
		Logger:    logger,
		Writer:    writer,
		flags:     InitFlags(),
		release:   release,
	}
}

func New(resources fs.FS, release Release) icli {
	return NewWithOptions(
		resources,
		log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.ParseLevel(os.Getenv("INITIUM_LOG_LEVEL")),
			ReportCaller:    true,
			ReportTimestamp: true,
		}),
		os.Stdout,
		release,
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
	projectType := cCtx.String(projectTypeFlag)
	projectDirectory := cCtx.String(projectDirectoryFlag)
	absProjectDirectory, err := filepath.Abs(cCtx.String(projectDirectoryFlag))
	registry := cCtx.String(repoNameFlag)
	imagePullSecrets := cCtx.StringSlice(imagePullSecretsFlag)
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
		project.ProjectType(projectType),
		projectDirectory,
		cCtx.String(runtimeVersionFlag),
		version,
		cCtx.Bool(isPrivateServiceFlag),
		imagePullSecrets,
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
	c.project = &project
	return nil
}

func (c *icli) getProject(cCtx *cli.Context) (*project.Project, error) {
	if c.project == nil {
		err := c.init(cCtx)
		if err != nil {
			return nil, err
		}
	}
	return c.project, nil
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
			c.SecretsCMD(),
			c.InitCMD(),
			{
				Name:  "version",
				Usage: "Return the version of the cli",
				Action: func(ctx *cli.Context) error {
					_, err := fmt.Fprintf(c.Writer, "version %s, commit %s, built at %s\n", c.release.Version, c.release.Commit, c.release.Date)
					return err
				},
			},
		},
		Before: func(ctx *cli.Context) error {
			if err := c.loadFlagsFromConfig(ctx); err != nil {
				return err
			}

			return c.checkRequiredFlags(ctx, []string{})
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	return app.Run(args)
}
