package cli

import (
	"fmt"
	"os"
	"sort"

	"github.com/charmbracelet/log"
	"k8s.io/utils/strings/slices"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"
	"github.com/urfave/cli/v2"
)

func (c CLI) InitGithubCMD(cCtx *cli.Context) error {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	options := project.InitOptions{
		PipelineType:      cCtx.Command.Name,
		DestinationFolder: cCtx.String(destinationFolderFlag),
		DefaultBranch:     cCtx.String(defaultBranchFlag),
		AppName:           cCtx.String(appNameFlag),
		Repository:        cCtx.String(repoNameFlag),
		ProjectDirectory:  cCtx.String(projectDirectoryFlag),
	}
	data, err := project.ProjectInit(options, c.Resources)

	if err != nil {
		return err
	}

	for _, v := range data {
		logger.Info(v)
	}

	return nil
}

func (c CLI) InitConfigCMD(ctx *cli.Context) error {
	excludedFlags := []string{
		"help",
		namespaceFlag,
		configFileFlag,
		projectDirectoryFlag,
		destinationFolderFlag,
		tokenFlag,
		registryPasswordFlag,
		caCRTFlag,
	}
	f := []cli.Flag{}
	for _, vs := range flags {
		f = append(f, vs...)
	}
	sort.Sort(cli.FlagsByName(f))

	for _, flag := range f {
		stringFlag := flag.(*cli.StringFlag)
		if slices.Contains(excludedFlags, stringFlag.Name) {
			continue
		}

		value := ctx.String(stringFlag.Name)
		if value == "" {
			stringFlag.GetValue()
		}

		if value == "" {
			fmt.Fprintf(c.Writer, "%s: null\n", stringFlag.Name)
		} else {
			fmt.Fprintf(c.Writer, "%s: %s\n", stringFlag.Name, ctx.String(stringFlag.Name))
		}
	}

	return nil
}

func (c CLI) InitServiceAccountCMD(ctx *cli.Context) error {
	return k8s.GetServiceAccount(c.Resources)
}

func (c CLI) InitCMD() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "create configuration for the cli [EXPERIMENTAL]",
		Subcommands: []*cli.Command{
			{
				Name:   "github",
				Usage:  "create a github pipeline yaml file",
				Flags:  c.CommandFlags(InitGithub),
				Action: c.InitGithubCMD,
				Before: c.baseBeforeFunc,
			},
			{
				Name:   "config",
				Usage:  "create a config file with all available flags set to null",
				Action: c.InitConfigCMD,
				Before: c.baseBeforeFunc,
			},
			{
				Name:   "service-account",
				Usage:  "output all resources needed to create a service account",
				Action: c.InitServiceAccountCMD,
				Before: c.baseBeforeFunc,
			},
		},
	}
}
