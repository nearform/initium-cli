package cli

import (
	"fmt"
	"github.com/nearform/initium-cli/src/services/versions"

	"github.com/nearform/initium-cli/src/services/git"
	"github.com/nearform/initium-cli/src/utils"
	"github.com/urfave/cli/v2"
)

const (
	cleanFlag      string = "clean"
	branchNameFlag string = "branch-name"
)

func (c icli) buildPushDeploy(cCtx *cli.Context) error {
	err := c.Build(cCtx)
	if err != nil {
		return fmt.Errorf("building %v", err)
	}
	if cCtx.Bool(stopOnBuildFlag) {
		return err
	}

	err = c.Push(cCtx)
	if err != nil {
		return fmt.Errorf("pushing %v", err)
	}
	if cCtx.Bool(stopOnPushFlag) {
		return err
	}
	return c.Deploy(cCtx)
}

func (c icli) OnBranchCMD() *cli.Command {
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
			Usage: "Stop the onbranch command after the build step",
		},
		&cli.BoolFlag{
			Name:  stopOnPushFlag,
			Value: false,
			Usage: "Stop the onbranch command after the push step",
		},
		&cli.BoolFlag{
			Name:  cleanFlag,
			Value: false,
			Usage: "Delete the knative service",
		},
		&cli.StringFlag{
			Name:  branchNameFlag,
			Usage: "Pass a branch name and disable autodetection",
		},
	}...)
	return &cli.Command{
		Name:  "onbranch",
		Usage: "Build, push and deploy the application using the branch name as version and namespace",
		Flags: flags,
		Action: func(cCtx *cli.Context) error {
			if cCtx.Bool(cleanFlag) {
				return c.Delete(cCtx)
			}

			return c.buildPushDeploy(cCtx)
		},
		Before: func(ctx *cli.Context) error {
			var err error
			if err := c.loadFlagsFromConfig(ctx); err != nil {
				return err
			}

			clientConfigFileName := ctx.String(configFileFlag)
			if err := versions.CheckClientConfigFileSchemaMatchesCli(clientConfigFileName, c.Resources); err != nil {
				return err
			}

			branchName := ctx.String(branchNameFlag)

			if branchName == "" {
				branchName, err = git.GetBranchName()
				if err != nil {
					return err
				}
			}

			ctx.Set(appVersionFlag, utils.EncodeRFC1123(branchName))
			ctx.Set(namespaceFlag, utils.EncodeRFC1123(branchName))

			ignoredFlags := []string{}
			if ctx.Bool(stopOnBuildFlag) {
				ignoredFlags = append(ignoredFlags, []string{registryPasswordFlag, registryUserFlag}...)
			}
			if ctx.Bool(stopOnPushFlag) {
				ignoredFlags = append(ignoredFlags, []string{endpointFlag, tokenFlag, caCRTFlag, namespaceFlag}...)
			}
			if ctx.Bool(cleanFlag) {
				ignoredFlags = append(ignoredFlags, []string{registryPasswordFlag, registryUserFlag}...)
			}

			return c.checkRequiredFlags(ctx, ignoredFlags)
		},
	}
}
