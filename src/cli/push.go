package cli

import (
	"github.com/docker/docker/api/types"
	"github.com/urfave/cli/v2"
)

func (c *icli) Push(cCtx *cli.Context) error {
	c.init(cCtx)
	c.DockerService.AuthConfig = types.AuthConfig{
		Username: cCtx.String(registryUserFlag),
		Password: cCtx.String(registryPasswordFlag),
	}
	return c.DockerService.Push()
}

func (c *icli) PushCMD() *cli.Command {
	return &cli.Command{
		Name:   "push",
		Usage:  "push the container image to a registry",
		Flags:  c.CommandFlags([]FlagsType{Registry, Shared}),
		Action: c.Push,
		Before: c.baseBeforeFunc,
	}
}
