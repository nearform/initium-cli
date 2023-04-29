package cli

import (
	knative "github.com/nearform/k8s-kurated-addons-cli/src/services/k8s"
	"github.com/urfave/cli/v2"
)

func (c CLI) Delete() *cli.Command {
	return &cli.Command{
		Name:  "delete",
		Usage: "delete the knative service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "endpoint",
				EnvVars: []string{"KKA_ENDPOINT"},
			},
			&cli.StringFlag{
				Name:    "token",
				EnvVars: []string{"KKA_TOKEN"},
			},
			&cli.StringFlag{
				Name:    "ca-crt",
				EnvVars: []string{"KKA_CA_CERT"},
			},
		},
		Action: func(cCtx *cli.Context) error {
			config, err := knative.Config(
				cCtx.String("endpoint"),
				cCtx.String("token"),
				[]byte(cCtx.String("ca-crt")),
			)
			if err != nil {
				return err
			}
			project := c.newProject(cCtx)
			return knative.Clean(config, project)
		},
	}
}
