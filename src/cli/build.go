package cli

import (
	"path"

	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"
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
	return &cli.Command{
		Name:   "build",
		Usage:  "build a container image from the project directory",
		Flags:  c.CommandFlags(Build),
		Action: c.Build,
		Before: c.baseBeforeFunc,
	}
}
