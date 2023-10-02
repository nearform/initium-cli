package cli

import (
	"path"

	"github.com/nearform/initium-cli/src/utils/logger"
	"github.com/urfave/cli/v2"
)

func (c *icli) Build(cCtx *cli.Context) error {
	project, err := c.getProject(cCtx)
	if err != nil {
		return err
	}

	logger.PrintInfo("Dockerfile Location: " + path.Join(project.Directory, c.DockerService.DockerFileName))
	return c.DockerService.Build()
}

func (c icli) BuildCMD() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "build a container image from the project directory",
		Flags:  c.CommandFlags([]FlagsType{Build, Shared}),
		Action: c.Build,
		Before: c.baseBeforeFunc,
	}
}
