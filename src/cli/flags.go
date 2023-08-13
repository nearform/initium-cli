package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nearform/initium-cli/src/services/git"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
	"k8s.io/utils/strings/slices"
)

type FlagsType string

const (
	Build      FlagsType = "build"
	Kubernetes FlagsType = "kubernetes"
	Registry   FlagsType = "registry"
	InitGithub FlagsType = "init-github"
	App        FlagsType = "app"
)

const (
	runtimeVersionFlag    string = "runtime-version"
	endpointFlag          string = "cluster-endpoint"
	tokenFlag             string = "cluster-token"
	caCRTFlag             string = "cluster-ca-crt"
	registryUserFlag      string = "registry-user"
	registryPasswordFlag  string = "registry-password"
	destinationFolderFlag string = "destination-folder"
	defaultBranchFlag     string = "default-branch"
	appNameFlag           string = "app-name"
	appVersionFlag        string = "app-version"
	projectDirectoryFlag  string = "project-directory"
	repoNameFlag          string = "repo-name"
	dockerFileNameFlag    string = "dockerfile-name"
	configFileFlag        string = "config-file"
	namespaceFlag         string = "namespace"
	stopOnBuildFlag       string = "stop-on-build"
	stopOnPushFlag        string = "stop-on-push"
)

var requiredFlags []string
var flags map[FlagsType]([]cli.Flag)

// This function is executed when the module is loaded
func init() {
	registry := ""
	org, err := git.GetGithubOrg()
	if err == nil {
		registry = fmt.Sprintf("ghcr.io/%s", org)
	}

	defaultFlags := map[FlagsType]([]cli.Flag){
		Build: []cli.Flag{
			&cli.StringFlag{
				Name:     runtimeVersionFlag,
				EnvVars:  []string{"INITIUM_RUNTIME_VERSION"},
				Category: "build",
			},
		},
		Kubernetes: []cli.Flag{
			&cli.StringFlag{
				Name:     endpointFlag,
				EnvVars:  []string{"INITIUM_CLUSTER_ENDPOINT"},
				Required: true,
				Category: "deploy",
			},
			&cli.StringFlag{
				Name:     tokenFlag,
				EnvVars:  []string{"INITIUM_CLUSTER_TOKEN"},
				Required: true,
				Category: "deploy",
			},
			&cli.StringFlag{
				Name:     caCRTFlag,
				EnvVars:  []string{"INITIUM_CLUSTER_CA_CERT"},
				Required: true,
				Category: "deploy",
			},
			&cli.StringFlag{
				Name:     namespaceFlag,
				EnvVars:  []string{"INITIUM_NAMESPACE"},
				Required: true,
				Category: "deploy",
			},
		},
		Registry: []cli.Flag{
			&cli.StringFlag{
				Name:     registryUserFlag,
				EnvVars:  []string{"INITIUM_REGISTRY_USER"},
				Required: true,
				Category: "registry",
			},
			&cli.StringFlag{
				Name:     registryPasswordFlag,
				EnvVars:  []string{"INITIUM_REGISTRY_PASSWORD"},
				Required: true,
				Category: "registry",
			},
		},
		InitGithub: []cli.Flag{
			&cli.StringFlag{
				Name:     destinationFolderFlag,
				Usage:    "Define a destination folder to place your pipeline file",
				Category: "init",
				Value:    defaults.GithubActionFolder,
			},
			&cli.StringFlag{
				Name:     defaultBranchFlag,
				Usage:    "Define the default branch in your repo",
				Category: "init",
				Value:    defaults.GithubDefaultBranch,
			},
		},
		App: []cli.Flag{
			&cli.StringFlag{
				Name:     appNameFlag,
				Usage:    "The name of the app",
				Required: true,
				EnvVars:  []string{"INITIUM_APP_NAME"},
			},
			&cli.StringFlag{
				Name:    appVersionFlag,
				Usage:   "The version of your application",
				Value:   defaults.AppVersion,
				EnvVars: []string{"INITIUM_VERSION"},
			},
			&cli.StringFlag{
				Name:    projectDirectoryFlag,
				Usage:   "The directory in which your Dockerfile lives",
				Value:   defaults.ProjectDirectory,
				EnvVars: []string{"INITIUM_PROJECT_DIRECTORY"},
			},
			&cli.StringFlag{
				Name:     repoNameFlag,
				Usage:    "The base address of the container repository",
				Value:    registry,
				Required: registry == "",
				EnvVars:  []string{"INITIUM_REPO_NAME"},
			},
			&cli.StringFlag{
				Name:    dockerFileNameFlag,
				Usage:   "The name of the Dockerfile",
				EnvVars: []string{"INITIUM_DOCKERFILE_NAME"},
			},
			&cli.StringFlag{
				Name:    configFileFlag,
				Usage:   "read parameters from config",
				Hidden:  true,
				Value:   defaults.ConfigFile,
				EnvVars: []string{"INITIUM_CONFIG_FILE"},
			},
		},
	}

	// TODO: urfave has an issue with required flags and altsrc https://github.com/urfave/cli/issues/1725
	// this is an hack to go around that issue, we should remove it once urfave v3 is released
	for _, vs := range defaultFlags {
		for _, flag := range vs {
			stringFlag := flag.(*cli.StringFlag)
			if stringFlag.Required {
				requiredFlags = append(requiredFlags, stringFlag.Name)
				stringFlag.Required = false
			}

		}
	}

	flags = defaultFlags
}

func (c CLI) checkRequiredFlags(ctx *cli.Context, ignoredFlags []string) error {
	missingFlags := []string{}

	for _, v := range ctx.Command.Flags {
		name := v.Names()[0]
		c.Logger.Debugf("%s is set to %s", name, ctx.String(name))
		if slices.Contains(requiredFlags, name) && !slices.Contains(ignoredFlags, name) && !ctx.IsSet(name) {
			missingFlags = append(missingFlags, name)
		}
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("required flags \"%v\" not set", strings.Join(missingFlags, ", "))
	}

	return nil
}

func (c CLI) loadFlagsFromConfig(ctx *cli.Context) error {
	config := make(map[interface{}]interface{})
	cfgFile := ctx.String(configFileFlag)
	//if the default config file doesn't exist we can ignore the rest and return nil
	_, err := os.Stat(cfgFile)
	if err != nil && errors.Is(err, os.ErrNotExist) && cfgFile == defaults.ConfigFile {
		return nil
	}

	yamlFile, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return fmt.Errorf("cannot read config file %s", cfgFile)
	}

	if err = yaml.Unmarshal(yamlFile, &config); err != nil {
		return fmt.Errorf("cannot parse config file %s, not a valid yaml", cfgFile)
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

func (c CLI) CommandFlags(command FlagsType) []cli.Flag {
	return flags[command]
}
