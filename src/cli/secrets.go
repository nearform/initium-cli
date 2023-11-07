package cli

import (
	"github.com/nearform/initium-cli/src/services/k8s"
	"github.com/nearform/initium-cli/src/utils/logger"
	"github.com/urfave/cli/v2"
	"k8s.io/client-go/rest"
)

func (c *icli) ListSecretsCMD(cCtx *cli.Context) error {
	config, err := setK8sConfig(cCtx)
	if err != nil {
		return err
	}

	stringSecretList, err := k8s.ListSecrets(
		config, cCtx.String(namespaceFlag))
	if err != nil {
		return err
	}

	logger.PrintInfo(stringSecretList)
	return nil
}

func (c *icli) CreateSecretsCMD(cCtx *cli.Context) error {
	config, err := setK8sConfig(cCtx)
	if err != nil {
		return err
	}

	// TODO: Secret name must be provided as argument
	// TODO: Validate at least 1 key/value
	// TODO: Validate same number of keys/values
	// TODO: Use all flag values
	err = k8s.CreateSecret(config, cCtx.String("name"), cCtx.StringSlice("key")[0], cCtx.StringSlice("value")[0], cCtx.String(namespaceFlag))
	if err != nil {
		return err
	}

	return nil
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
			{
				Name:   "create",
				Usage:  "create a secret",
				Action: c.CreateSecretsCMD,
				Before: c.baseBeforeFunc,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "name",
						Usage: "Secret name",
					},					
					&cli.StringSliceFlag{
						Name:  "key",
						Usage: "Secret key. Must be used in conjunction with --value. Allows multiple, and order matters",
					},
					&cli.StringSliceFlag{
						Name:  "value",
						Usage: "Secret value. Allows multiple, and they must be in the same order as the respective --key",
					},
				},
			},
		},
	}
}

func setK8sConfig(cCtx *cli.Context) (*rest.Config, error) {
	config, err := k8s.Config(
		cCtx.String(endpointFlag),
		cCtx.String(tokenFlag),
		[]byte(cCtx.String(caCRTFlag)),
	)
	return config, err
}
