package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *CLI) OnMainCMD() *cli.Command {
	flags := []cli.Flag{}
	flags = append(flags, c.CommandFlags(Kubernetes)...)
	flags = append(flags, c.CommandFlags(Build)...)
	flags = append(flags, c.CommandFlags(Registry)...)
	return &cli.Command{
		Name:  "onmain",
		Usage: "deploy the application as a knative service",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			err := c.Build(cCtx)
			if err != nil {
				return fmt.Errorf("Building %v", err)
			}

			err = c.Push(cCtx)
			if err != nil {
				return fmt.Errorf("Pushing %v", err)
			}

			return c.Deploy(cCtx)
		},
		Before: func(ctx *cli.Context) error {
			err := c.loadFlagsFromConfig(ctx)

			if err != nil {
				c.Logger.Debug("failed to load config", err)
			}

			return nil
		},
	}
}
