package cli

import (
	"bufio"
	"fmt"
	"github.com/nearform/initium-cli/src/services/git"
	knative "github.com/nearform/initium-cli/src/services/k8s"
	"github.com/nearform/initium-cli/src/utils"
	"github.com/urfave/cli/v2"
	"os"
)

func (c *icli) AppSetSecret(cCtx *cli.Context) error {
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

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please provide a value for the secret key: ")
	secretValue, _ := reader.ReadString('\n')

	// Remove newline characters from the input
	secretValue = secretValue[:len(secretValue)-1]

	return knative.SecretUpd(cCtx.String(secretKeyFlag), secretValue, config, project, cCtx.String(namespaceFlag))
}

func (c *icli) AppSecretCMD() *cli.Command {
	return &cli.Command{
		Name:   "appsecret",
		Usage:  "creates K8s secret & assigns it to knative app",
		Flags:  c.CommandFlags([]FlagsType{Kubernetes, Shared, AppSecret}),
		Action: c.AppSetSecret,
		Before: func(ctx *cli.Context) error {
			if err := c.loadFlagsFromConfig(ctx); err != nil {
				return err
			}

			branchName, err := git.GetBranchName()
			if err != nil {
				return err
			}

			ctx.Set(namespaceFlag, utils.EncodeRFC1123(branchName))

			if err := c.checkRequiredFlags(ctx, []string{}); err != nil {
				return err
			}
			return nil
		},
	}
}
