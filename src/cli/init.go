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

func (c CLI) Init(cCtx *cli.Context) error {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)
	options := project.InitOptions{
		PipelineType:      cCtx.Command.Name,
		DestinationFolder: cCtx.String("destination-folder"),
		DefaultBranch:     cCtx.String("default-branch"),
		AppName:           cCtx.String("app-name"),
		Repository:        cCtx.String("repo-name"),
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

func (c CLI) InitConfig(ctx *cli.Context) error {
	excludedFlags := []string{
		"help",
		"config-file",
		"project-directory",
		"destination-folder",
		"token",
		"registry-password",
		"ca-crt",
	}
	flags := []cli.Flag{}
	for _, vs := range c.Flags() {
		flags = append(flags, vs...)
	}
	sort.Sort(cli.FlagsByName(flags))

	for _, flag := range flags {
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

func (c CLI) InitCMD() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "create configuration for the cli [EXPERIMENTAL]",
		Subcommands: []*cli.Command{
			{
				Name:   "github",
				Usage:  "create a github pipeline yaml file",
				Flags:  c.CommandFlags(InitGithub),
				Action: c.Init,
				Before: func(ctx *cli.Context) error {
					err := c.loadFlagsFromConfig(ctx)

					if err != nil {
						c.Logger.Debug("failed to load config", err)
					}

					return nil
				},
			},
			{
				Name:   "config",
				Usage:  "create a config file with all available flags set to null",
				Action: c.InitConfig,
				Before: func(ctx *cli.Context) error {
					err := c.loadFlagsFromConfig(ctx)

					if err != nil {
						c.Logger.Debug("failed to load config", err)
					}

					return nil
				},
			},
		},
	}
}
