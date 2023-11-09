package cli

import (
	"fmt"

	"github.com/nearform/initium-cli/src/services/secrets"
	"github.com/urfave/cli/v2"
)

func (c icli) generateKeys(ctx *cli.Context) error {
	keys, err := secrets.GenerateKeys()
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Writer, "Private key: %q\n", keys.Private)
	fmt.Fprintf(c.Writer, "Public key: %q\n", keys.Public)
	return nil
}

func (c icli) encrypt(ctx *cli.Context) error {
	publicKey := ctx.String("publicKey")
	secretMaterial := ctx.String("secret")
	return secrets.Encrypt(publicKey, secretMaterial, c.Writer)
}

func (c icli) decrypt(ctx *cli.Context) error {
	privateKey := ctx.String("privateKey")
	secret := ctx.String("secret")
	return secrets.Decrypt(privateKey, secret, c.Writer)
}

func (c icli) SecretsCMD() *cli.Command {

	return &cli.Command{
		Name:  "secrets",
		Usage: "create configuration for the cli [EXPERIMENTAL]",
		Subcommands: []*cli.Command{
			{
				Name:   "generate-keys",
				Usage:  "generate the public and private keys",
				Action: c.generateKeys,
				Before: c.baseBeforeFunc,
			},
			{
				Name:   "encrypt",
				Action: c.encrypt,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "publicKey",
						EnvVars: []string{"INITIUM_SECRET_PUBLIC_KEY"},
					},
					&cli.StringFlag{
						Name:    "secret",
						EnvVars: []string{"INITIUM_SECRET_MATERIAL"},
					},
				},
			},
			{
				Name:   "decrypt",
				Action: c.decrypt,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "privateKey",
						EnvVars: []string{"INITIUM_SECRET_PRIVATE_KEY"},
					},
					&cli.StringFlag{
						Name:    "secret",
						EnvVars: []string{"INITIUM_ENCRYPTED_SECRET"},
					},
				},
			},
		},
	}
}
