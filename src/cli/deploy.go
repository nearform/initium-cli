package cli

import (
	"github.com/nearform/initium-cli/src/services/git"
	knative "github.com/nearform/initium-cli/src/services/k8s"
	"github.com/urfave/cli/v2"
)

func (c *icli) Deploy(cCtx *cli.Context) error {
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

	commitSha, err := git.GetHash()
	if err != nil {
		return err
	}

	return knative.Apply(cCtx.String(namespaceFlag), commitSha, config, project, c.dockerImage)
}

func (c icli) DeployCMD() *cli.Command {
	return &cli.Command{
		Name:   "deploy",
		Usage:  "deploy the application as a knative service",
		Flags:  c.CommandFlags([]FlagsType{Kubernetes, Shared}),
		Action: c.Deploy,
		Before: c.baseBeforeFunc,
	}
}
