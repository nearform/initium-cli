package cli

import (
	"github.com/docker/docker/api/types"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	"github.com/urfave/cli/v2"
)

func (c CLI) Push() *cli.Command {
	return &cli.Command{
		Name:  "push",
		Usage: "push the container image to a registry",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "registry-user",
				EnvVars: []string{"KKA_REGISTRY_USER"},
			},
			&cli.StringFlag{
				Name:    "registry-password",
				EnvVars: []string{"KKA_REGISTRY_PASSWORD"},
			},
		},
		Action: func(cCtx *cli.Context) error {
			repoName := cCtx.String("repo-name")
			dockerFileName := cCtx.String("dockerfile-name")
			project := c.newProject(cCtx)
			docker := docker.New(project, dockerFileName, repoName)
			docker.AuthConfig = types.AuthConfig{
				Username: cCtx.String("registry-user"),
				Password: cCtx.String("registry-password"),
			}
			return docker.Push()
		},
	}
}
