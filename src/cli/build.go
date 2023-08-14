package cli

import (
	"path"

	"github.com/nearform/initium-cli/src/utils/logger"
	"github.com/urfave/cli/v2"
)

func (c *CLI) Build(cCtx *cli.Context) error {
	project, err := c.getProject(cCtx)
	if err != nil {
		return err
	}

	logger.PrintInfo("Dockerfile Location: " + path.Join(project.Directory, c.DockerService.DockerFileName))
	return c.DockerService.Build()
}

func (c CLI) BuildCMD() *cli.Command {
	flags := []cli.Flag{}
	flags = append(flags, c.CommandFlags(Build)...)
	flags = append(flags, c.CommandFlags(Shared)...)

	return &cli.Command{
		Name:   "build",
		Usage:  "build a container image from the project directory",
		Flags:  flags,
		Action: c.Build,
		Before: c.baseBeforeFunc,
	}
}
