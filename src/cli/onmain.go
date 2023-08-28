package cli

import (
	"github.com/nearform/initium-cli/src/utils/defaults"
	"github.com/urfave/cli/v2"
)

func (c *icli) OnMainCMD() *cli.Command {
	flags := c.CommandFlags([]FlagsType{
		Kubernetes,
		Build,
		Registry,
		Shared,
	})
	flags = append(flags, []cli.Flag{
		&cli.BoolFlag{
			Name:  stopOnBuildFlag,
			Value: false,
		},
		&cli.BoolFlag{
			Name:  stopOnPushFlag,
			Value: false,
		},
	}...)

	return &cli.Command{
		Name:   "onmain",
		Usage:  "deploy the application as a knative service",
		Flags:  flags,
		Action: c.buildPushDeploy,
		Before: func(ctx *cli.Context) error {
			if err := c.loadFlagsFromConfig(ctx); err != nil {
				return err
			}

			ctx.Set(appVersionFlag, "latest")
			ctx.Set(namespaceFlag, defaults.GithubDefaultBranch)

			ignoredFlags := []string{}
			if ctx.Bool(stopOnBuildFlag) {
				ignoredFlags = append(ignoredFlags, []string{registryPasswordFlag, registryUserFlag}...)
			}
			if ctx.Bool(stopOnPushFlag) {
				ignoredFlags = append(ignoredFlags, []string{endpointFlag, tokenFlag, caCRTFlag, namespaceFlag}...)
			}

			return c.checkRequiredFlags(ctx, ignoredFlags)
		},
	}
}
