package cli

import (
	knative "github.com/nearform/initium-cli/src/services/k8s"
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
	project, err := c.getProject(cCtx)
	if err != nil {
		return err
	}
	return knative.Clean(cCtx.String(namespaceFlag), config, project)
}

func (c *CLI) DeleteCMD() *cli.Command {
	return &cli.Command{
		Name:   "delete",
		Usage:  "delete the knative service",
		Flags:  c.CommandFlags([]FlagsType{Kubernetes, Shared}),
		Action: c.Delete,
		Before: c.baseBeforeFunc,
	}
}
