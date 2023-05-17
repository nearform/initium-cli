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
	project := c.newProject(cCtx)
	return knative.Clean(config, project)
}

func (c *CLI) DeleteCMD() *cli.Command {
	return &cli.Command{
		Name:   "delete",
		Usage:  "delete the knative service",
		Flags:  Flags(Kubernetes),
		Action: c.Delete,
	}
}
