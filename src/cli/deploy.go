package cli

import (
	knative "github.com/nearform/initium-cli/src/services/k8s"
	"github.com/urfave/cli/v2"
)

func (c *CLI) Deploy(cCtx *cli.Context) error {
	config, err := knative.Config(
		cCtx.String(endpointFlag),
		cCtx.String(tokenFlag),
		[]byte(cCtx.String(caCRTFlag)),
	)

	if err != nil {
		return err
	}
	project, err := c.getProject(cCtx)
	if err != nil {
		return err
	}

	return knative.Apply(cCtx.String(namespaceFlag), config, project, c.dockerImage)
}

func (c CLI) DeployCMD() *cli.Command {
	return &cli.Command{
		Name:   "deploy",
		Usage:  "deploy the application as a knative service",
		Flags:  c.CommandFlags(Kubernetes),
		Action: c.Deploy,
		Before: c.baseBeforeFunc,
	}
}
