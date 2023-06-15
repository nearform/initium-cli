package cli

import (
	knative "github.com/nearform/k8s-kurated-addons-cli/src/services/k8s"
	"github.com/urfave/cli/v2"
)

func (c *CLI) Delete(cCtx *cli.Context) error {
	config, err := knative.Config(
		cCtx.String("endpoint"),
		cCtx.String("token"),
		[]byte(cCtx.String("ca-crt")),
	)
	if err != nil {
		return err
	}
	project := c.getProject(cCtx)
	return knative.Clean(config, project)
}

func (c *CLI) DeleteCMD() *cli.Command {
	return &cli.Command{
		Name:   "delete",
		Usage:  "delete the knative service",
		Flags:  c.CommandFlags(Kubernetes),
		Action: c.Delete,
		Before: func(ctx *cli.Context) error {
			err := c.loadFlagsFromConfig(ctx)

			if err != nil {
				c.Logger.Debug("failed to load config", err)
			}

			return nil
		},
	}
}
