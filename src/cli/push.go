package cli

import (
	"github.com/docker/docker/api/types"
	"github.com/urfave/cli/v2"
)

func (c *CLI) Push(cCtx *cli.Context) error {
	c.init(cCtx)
	c.DockerService.AuthConfig = types.AuthConfig{
		Username: cCtx.String(registryUserFlag),
		Password: cCtx.String(registryPasswordFlag),
	}
	return c.DockerService.Push()
}

func (c *CLI) PushCMD() *cli.Command {
	flags := []cli.Flag{}
	flags = append(flags, c.CommandFlags(Registry)...)
	flags = append(flags, c.CommandFlags(Shared)...)

	return &cli.Command{
		Name:   "push",
		Usage:  "push the container image to a registry",
		Flags:  flags,
		Action: c.Push,
		Before: c.baseBeforeFunc,
	}
}
