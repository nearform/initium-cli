package cli

import (
	knative "github.com/nearform/k8s-kurated-addons-cli/src/services/k8s"
	"github.com/urfave/cli/v2"
)

func (c *CLI) Delete(cCtx *cli.Context) error {
	config, err := knative.Config(
		cCtx.String(endpointFlag),
		cCtx.String(tokenFlag),
		[]byte(cCtx.String(caCRTFlag)),
	)
	if err != nil {
		return err
	}
	project := c.getProject(cCtx)
	return knative.Clean(cCtx.String(namespaceFlag), config, project)
}

func (c *CLI) DeleteCMD() *cli.Command {
	return &cli.Command{
		Name:   "delete",
		Usage:  "delete the knative service",
		Flags:  c.CommandFlags(Kubernetes),
		Action: c.Delete,
		Before: c.baseBeforeFunc,
	}
}
