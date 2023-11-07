package cli

import (
	"github.com/nearform/initium-cli/src/services/k8s"
	"github.com/urfave/cli/v2"
)

func (c *icli) ListSecretsCMD(cCtx *cli.Context) error {

	config, err := k8s.Config(
		cCtx.String(endpointFlag),
		cCtx.String(tokenFlag),
		[]byte(cCtx.String(caCRTFlag)),
	)
	if err != nil {
		return err
	}

	_, err = k8s.ListSecrets(
		config, cCtx.String(namespaceFlag))
	if err != nil {
		return err
	}

	return nil
	// Convert to string separating each item on the proper line
	// Convert secrets to a string that will show a result similar to `kubectl get secrets`

}

func (c icli) SecretsCMD() *cli.Command {
	flags := c.CommandFlags([]FlagsType{
		Kubernetes,
		Shared,
	})

	flags = append(flags, []cli.Flag{
		&cli.StringFlag{
			Name:  branchNameFlag,
			Usage: "Pass a branch name and disable autodetection",
		},
	}...)

	return &cli.Command{
		Name:  "secrets",
		Usage: "manage k8s secrets for apps deployed through Initium", // TODO: Consider a 'initium' prefix for secrets created through this command
		Flags: flags,
		Before: func(ctx *cli.Context) error {
			var err error
			err = c.detectFlagsFromGit(ctx)
			if err != nil {
				return err
			}

			err = c.baseBeforeFunc(ctx)
			if err != nil {
				return err
			}
			return nil
		},
		Subcommands: []*cli.Command{
			{
				Name:   "list",
				Usage:  "list secrets managed by initium",
				Action: c.ListSecretsCMD,
				Before: c.baseBeforeFunc,
			},
		},
	}
}
