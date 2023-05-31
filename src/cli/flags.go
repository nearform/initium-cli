package cli

import (
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
	"github.com/urfave/cli/v2"
)

type FlagsType string

const (
	Build      FlagsType = "build"
	Kubernetes FlagsType = "kubernetes"
	Registry   FlagsType = "registry"
	InitGithub FlagsType = "init-github"
)

func Flags(command FlagsType) []cli.Flag {
	flags := map[FlagsType]([]cli.Flag){
		Build: []cli.Flag{
			&cli.StringFlag{
				Name:     "runtime-version",
				EnvVars:  []string{"KKA_RUNTIME_VERSION"},
				Category: "build",
			},
		},
		Kubernetes: []cli.Flag{
			&cli.StringFlag{
				Name:     "endpoint",
				EnvVars:  []string{"KKA_ENDPOINT"},
				Required: true,
				Category: "deploy",
			},
			&cli.StringFlag{
				Name:     "token",
				EnvVars:  []string{"KKA_TOKEN"},
				Required: true,
				Category: "deploy",
			},
			&cli.StringFlag{
				Name:     "ca-crt",
				EnvVars:  []string{"KKA_CA_CERT"},
				Required: true,
				Category: "deploy",
			},
		},
		Registry: []cli.Flag{
			&cli.StringFlag{
				Name:     "registry-user",
				EnvVars:  []string{"KKA_REGISTRY_USER"},
				Required: true,
				Category: "registry",
			},
			&cli.StringFlag{
				Name:     "registry-password",
				EnvVars:  []string{"KKA_REGISTRY_PASSWORD"},
				Required: true,
				Category: "registry",
			},
		},
		InitGithub: []cli.Flag{
			&cli.StringFlag{
				Name:     "destination-folder",
				Usage:    "Define a destination folder to place your pipeline file",
				Category: "init",
				Value:    defaults.GithubActionFolder,
			},
			&cli.StringFlag{
				Name:     "default-branch",
				Usage:    "Define the default branch in your repo",
				Category: "init",
				Value:    defaults.GithubDefaultBranch,
			},
		},
	}

	return flags[command]
}
