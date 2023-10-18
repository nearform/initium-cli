package cli

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nearform/initium-cli/src/services/git"
	"github.com/nearform/initium-cli/src/services/project"
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
	Shared     FlagsType = "shared"
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
	projectTypeFlag       string = "project-type"
	repoNameFlag          string = "container-registry"
	dockerFileNameFlag    string = "dockerfile-name"
	configFileFlag        string = "config-file"
	namespaceFlag         string = "namespace"
	imagePullSecretsFlag  string = "image-pull-secrets"
	stopOnBuildFlag       string = "stop-on-build"
	stopOnPushFlag        string = "stop-on-push"
	envVarFileFlag        string = "env-var-file"
)

type flags struct {
	requiredFlags []string
	all           map[FlagsType]([]cli.Flag)
}

func InitFlags() flags {
	registry := ""
	org, err := git.GetGithubOrg()
	if err == nil {
		registry = fmt.Sprintf("ghcr.io/%s", org)
	}

	appName := ""
	guess := project.GuessAppName()

	if guess != nil {
		appName = *guess
	}

	var projectType project.ProjectType
	tempProjectType, err := project.DetectType(".")
	if err == nil {
		projectType = tempProjectType
	}

	f := flags{
		requiredFlags: []string{},
		all: map[FlagsType]([]cli.Flag){
			Build: []cli.Flag{
				&cli.StringFlag{
					Name:     runtimeVersionFlag,
					EnvVars:  []string{"INITIUM_RUNTIME_VERSION"},
					Category: "build",
				},
				&cli.StringFlag{
					Name:     projectTypeFlag,
					Usage:    "The project type (node, go)",
					Value:    string(projectType),
					EnvVars:  []string{"INITIUM_PROJECT_TYPE"},
					Category: "build",
					Required: projectType == "",
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
				&cli.StringSliceFlag{
					Name:     imagePullSecretsFlag,
					Usage:    "Define one or more (repeating the flag or in csv format for the environment variable) image pull secrets",
					EnvVars:  []string{"INITIUM_IMAGE_PULL_SECRETS"},
					Required: false,
					Category: "deploy",
				},
				&cli.StringFlag{
					Name:     envVarFileFlag,
					Value:    defaults.EnvVarFile,
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
					Name:    projectDirectoryFlag,
					Usage:   "The directory in which your Dockerfile lives",
					Value:   defaults.ProjectDirectory,
					EnvVars: []string{"INITIUM_PROJECT_DIRECTORY"},
				},
				&cli.StringFlag{
					Name:    configFileFlag,
					Usage:   "read parameters from config",
					Value:   defaults.ConfigFile,
					EnvVars: []string{"INITIUM_CONFIG_FILE"},
				},
			},
			Shared: []cli.Flag{
				&cli.StringFlag{
					Name:     appNameFlag,
					Usage:    "The name of the app",
					Value:    appName,
					Required: appName == "",
					EnvVars:  []string{"INITIUM_APP_NAME"},
				},
				&cli.StringFlag{
					Name:    appVersionFlag,
					Usage:   "The version of your application",
					Value:   defaults.AppVersion,
					EnvVars: []string{"INITIUM_VERSION"},
				},
				&cli.StringFlag{
					Name:     repoNameFlag,
					Aliases:  []string{"repo-name"}, // keep compatibility with old version of the config
					Usage:    "The base address of the container registry",
					Value:    registry,
					Required: registry == "",
					EnvVars:  []string{"INITIUM_CONTAINER_REGISTRY", "INITIUM_REPO_NAME"}, // INITIUM_REPO_NAME to keep compatibility with older config
				},
				&cli.StringFlag{
					Name:    dockerFileNameFlag,
					Usage:   "The name of the Dockerfile",
					EnvVars: []string{"INITIUM_DOCKERFILE_NAME"},
				},
			},
		},
	}

	// TODO: urfave has an issue with required flags and altsrc https://github.com/urfave/cli/issues/1725
	// this is an hack to go around that issue, we should remove it once urfave v3 is released
	for _, vs := range f.all {
		for _, flag := range vs {
			switch flag.(type) {
			case *cli.StringFlag:
				stringFlag := flag.(*cli.StringFlag)
				if stringFlag.Required {
					f.requiredFlags = append(f.requiredFlags, stringFlag.Name)
					stringFlag.Required = false
				}
			}
		}
	}

	return f
}

func (c icli) checkRequiredFlags(ctx *cli.Context, ignoredFlags []string) error {
	missingFlags := []string{}
	requiredFlags := c.flags.requiredFlags

	for _, v := range ctx.Command.Flags {
		name := v.Names()[0]
		if slices.Contains(requiredFlags, name) && !slices.Contains(ignoredFlags, name) && !ctx.IsSet(name) {
			missingFlags = append(missingFlags, name)
		}
	}

	if len(missingFlags) > 0 {
		return fmt.Errorf("required flags \"%v\" not set", strings.Join(missingFlags, ", "))
	}

	return nil
}

func (c icli) loadFlagsFromConfig(ctx *cli.Context) error {
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

	excludedFlags := excludedFlagsFromConfig()

	c.Logger.Debugf("Loading flags %v", ctx.Command)
	for _, v := range ctx.Command.Flags {
		mainName := v.Names()[0]
		for _, name := range v.Names() {
			c.Logger.Debugf("%s is set = %v", name, ctx.IsSet(name))
			if name != "help" && !slices.Contains(excludedFlags, name) && config[name] != nil && !ctx.IsSet(mainName) {
				if err := ctx.Set(mainName, config[name].(string)); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (c icli) CommandFlags(commands []FlagsType) []cli.Flag {
	f := c.flags.all
	result := []cli.Flag{}

	for _, command := range commands {
		result = append(result, f[command]...)
	}

	return result
}
