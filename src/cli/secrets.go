package cli

import (
	"encoding/base64"
	"fmt"

	"github.com/nearform/initium-cli/src/services/secrets"
	"github.com/urfave/cli/v2"
)

func (c icli) generateKeys(ctx *cli.Context) error {
	keys, err := secrets.GenerateKeys()
	if err != nil {
		return err
	}

	fmt.Fprintf(c.Writer, "Secret key: %q\n", keys.Private)
	fmt.Fprintf(c.Writer, "Public key: %q\n", keys.Public)
	return nil
}

func (c icli) encrypt(ctx *cli.Context) error {
	publicKey := ctx.String(publicKeyFlag)
	secret := ctx.String(plainSecretFlag)
	base64Secret := ctx.String(base64PlainSecretFlag)

	if base64Secret == "" {
		base64Secret = base64.StdEncoding.EncodeToString([]byte(secret))
	}

	result, err := secrets.Encrypt(publicKey, base64Secret)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.Writer, "%s\n", result)
	return nil
}

func (c icli) decrypt(ctx *cli.Context) error {
	privateKey := ctx.String(privateKeyFlag)
	secret := ctx.String(base64EncryptedSecretFlag)
	result, err := secrets.Decrypt(privateKey, secret)
	if err != nil {
		return err
	}
	fmt.Fprintf(c.Writer, "%s\n", result)
	return nil
}

func (c icli) SecretsCMD() *cli.Command {

	return &cli.Command{
		Name:  "secrets",
		Usage: "A series of command to generate age keys, encrypt and decrypt secrets",
		Subcommands: []*cli.Command{
			{
				Name:   "generate-keys",
				Usage:  "Generate the public and private keys and output them on stdout",
				Action: c.generateKeys,
				Before: c.baseBeforeFunc,
			},
			{
				Name:   "encrypt",
				Usage:  "Encrypt a secret, if the secret flag is used the secret is first encoded in base64 and then encrypted",
				Action: c.encrypt,
				Flags:  c.CommandFlags([]FlagsType{Encrypt}),
				Before: func(ctx *cli.Context) error {
					if err := c.loadFlagsFromConfig(ctx); err != nil {
						return err
					}

					ignoredFlags := []string{}

					if ctx.IsSet(plainSecretFlag) {
						ignoredFlags = append(ignoredFlags, base64PlainSecretFlag)
					}
					if ctx.IsSet(base64PlainSecretFlag) {
						ignoredFlags = append(ignoredFlags, plainSecretFlag)
					}

					return c.checkRequiredFlags(ctx, ignoredFlags)
				},
			},
			{
				Name:   "decrypt",
				Usage:  "Decrypt a base64 encoded secret and output the base64 encoded value",
				Action: c.decrypt,
				Flags:  c.CommandFlags([]FlagsType{Decrypt}),
				Before: c.baseBeforeFunc,
			},
		},
	}
}
