package cli

import (
	knative "github.com/nearform/k8s-kurated-addons-cli/src/services/k8s"
	"github.com/urfave/cli/v2"
)

func (c *CLI) Deploy(cCtx *cli.Context) error {
	config, err := knative.Config(
		cCtx.String("endpoint"),
		cCtx.String("token"),
		[]byte(cCtx.String("ca-crt")),
	)

	if err != nil {
		return err
	}
	project := c.getProject(cCtx)

	return knative.Apply(config, project)
}

func (c CLI) DeployCMD() *cli.Command {
	return &cli.Command{
		Name:   "deploy",
		Usage:  "deploy the application as a knative service",
		Flags:  Flags(Kubernetes),
		Action: c.Deploy,
	}
}
