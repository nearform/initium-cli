package cli

import (
	"io/ioutil"
	"path"

	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type FlagsType string

const (
	Build      FlagsType = "build"
	Kubernetes FlagsType = "kubernetes"
	Registry   FlagsType = "registry"
	InitGithub FlagsType = "init-github"
	App        FlagsType = "app"
)

func (c CLI) loadFlagsFromConfig(ctx *cli.Context) error {
	config := make(map[interface{}]interface{})
	cfgFile := ctx.String("config-file")
	yamlFile, err := ioutil.ReadFile(cfgFile)

	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		return err
	}

	for _, v := range ctx.Command.Flags {
		name := v.Names()[0]
		c.Logger.Debugf("%s is set to %s", name, ctx.String(name))
		if name != "help" && !ctx.IsSet(name) {
			if config[name] != nil {
				c.Logger.Debugf("Loading %s as %s", name, config[name])
				ctx.Set(name, config[name].(string))
			}
		}
	}

	return nil
}

func (c CLI) Flags() map[FlagsType]([]cli.Flag) {
	return map[FlagsType]([]cli.Flag){
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
		App: []cli.Flag{
			&cli.StringFlag{
				Name:    "app-name",
				Usage:   "The name of the app",
				Value:   path.Base(c.CWD),
				EnvVars: []string{"KKA_APP_NAME"},
			},
			&cli.StringFlag{
				Name:    "app-version",
				Usage:   "The version of your application",
				Value:   "latest",
				EnvVars: []string{"KKA_VERSION"},
			},
			&cli.StringFlag{
				Name:    "project-directory",
				Usage:   "The directory in which your Dockerfile lives",
				Value:   defaults.ProjectDirectory,
				EnvVars: []string{"KKA_PROJECT_DIRECTORY"},
			},
			&cli.StringFlag{
				Name:    "repo-name",
				Usage:   "The base address of the container repository",
				Value:   defaults.RepoName,
				EnvVars: []string{"KKA_REPO_NAME"},
			},
			&cli.StringFlag{
				Name:    "dockerfile-name",
				Usage:   "The name of the Dockerfile",
				Value:   defaults.DockerfileName,
				EnvVars: []string{"KKA_DOCKERFILE_NAME"},
			},
			&cli.StringFlag{
				Name:    "config-file",
				Usage:   "read parameters from config",
				Value:   defaults.ConfigFile,
				EnvVars: []string{"KKA_CONFIG_FILE"},
			},
		},
	}
}

func (c CLI) CommandFlags(command FlagsType) []cli.Flag {
	return c.Flags()[command]
}
