package cli

import (
	"fmt"

	"github.com/nearform/initium-cli/src/services/git"
	knative "github.com/nearform/initium-cli/src/services/k8s"
	"github.com/urfave/cli/v2"
)

func (c *icli) Deploy(cCtx *cli.Context) error {
	namespace := cCtx.String(namespaceFlag)
	envFile := cCtx.String(envVarFileFlag)
	SecretRefEnvFile := cCtx.String(secretRefEnvFileFlag)

	project, err := c.getProject(cCtx)
	if err != nil {
		return err
	}

	commitSha, err := git.GetHash()
	if err != nil {
		return err
	}

	serviceManifest, err := knative.LoadManifest(namespace, commitSha, project, c.dockerImage, envFile, SecretRefEnvFile)
	if err != nil {
		return err
	}

	if cCtx.Bool(dryRunFlag) {
		yamlBytes, err := knative.ToYaml(serviceManifest)
		if err != nil {
			return err
		}
		fmt.Fprintf(c.Writer, "%s", yamlBytes)
		return nil
	}

	config, err := knative.Config(
		cCtx.String(endpointFlag),
		cCtx.String(tokenFlag),
		[]byte(cCtx.String(caCRTFlag)),
	)

	if err != nil {
		return err
	}

	err = knative.Apply(serviceManifest, config)
	return err
}

func (c icli) DeployCMD() *cli.Command {
	flags := c.CommandFlags([]FlagsType{Kubernetes, Shared})

	flags = append(flags, &cli.BoolFlag{
		Name:  dryRunFlag,
		Usage: "print out the knative manifest without applying it",
		Value: false,
	})

	return &cli.Command{
		Name:   "deploy",
		Usage:  "deploy the application as a knative service",
		Flags:  flags,
		Action: c.Deploy,
		Before: func(ctx *cli.Context) error {
			if err := c.loadFlagsFromConfig(ctx); err != nil {
				return err
			}

			ignoredFlags := []string{}
			if ctx.Bool(dryRunFlag) {
				ignoredFlags = append(ignoredFlags, endpointFlag, tokenFlag, caCRTFlag)
			}

			return c.checkRequiredFlags(ctx, ignoredFlags)
		},
	}
}
