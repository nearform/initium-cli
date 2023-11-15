package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/log"
	"k8s.io/utils/strings/slices"

	"github.com/nearform/initium-cli/src/services/k8s"
	"github.com/nearform/initium-cli/src/services/project"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"github.com/urfave/cli/v2"
)

const (
	persistFlag = "persist"
)

func (c icli) InitGithubCMD(cCtx *cli.Context) error {
	logger := log.New(os.Stderr)
	logger.SetLevel(log.DebugLevel)

	options := project.InitOptions{
		PipelineType:      cCtx.Command.Name,
		DestinationFolder: cCtx.String(destinationFolderFlag),
		DefaultBranch:     cCtx.String(defaultBranchFlag),
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

func excludedFlagsFromConfig() []string {
	return []string{
		"help",
		appVersionFlag,
		namespaceFlag,
		configFileFlag,
		projectDirectoryFlag,
		destinationFolderFlag,
		tokenFlag,
		registryPasswordFlag,
		caCRTFlag,
		registryUserFlag,
		endpointFlag,
	}
}

func (c icli) InitConfigCMD(ctx *cli.Context) error {
	excludedFlags := excludedFlagsFromConfig()

	f := []cli.Flag{}
	for _, vs := range c.flags.all {
		f = append(f, vs...)
	}

	sort.Sort(cli.FlagsByName(f))
	var n, v string
	config := ""
	for _, flag := range f {
		switch flag.(type) {
		case *cli.StringFlag:
			stringFlag := flag.(*cli.StringFlag)
			if slices.Contains(excludedFlags, stringFlag.Name) {
				continue
			}

			n = stringFlag.Name
			v = ctx.String(stringFlag.Name)
		case *cli.StringSliceFlag:
			stringSliceFlag := flag.(*cli.StringSliceFlag)
			if slices.Contains(excludedFlags, stringSliceFlag.Name) {
				continue
			}
			n = stringSliceFlag.Name
			v = strings.Join(ctx.StringSlice(stringSliceFlag.Name), ",")
		case *cli.BoolFlag:
			boolFlag := flag.(*cli.BoolFlag)
			if slices.Contains(excludedFlags, boolFlag.Name) {
				continue
			}
			n = boolFlag.Name
			v = strconv.FormatBool(ctx.Bool(boolFlag.Name))
		}

		if v == "" {
			v = "null"
		}
		config = config + fmt.Sprintf("%s: %s\n", n, v)
	}

	if ctx.Bool(persistFlag) {
		f, err := os.OpenFile(filepath.Join(ctx.String(projectDirectoryFlag), defaults.ConfigFile), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		f.WriteString(config)
	} else {
		fmt.Fprint(c.Writer, config)
	}

	return nil
}

func (c icli) InitServiceAccountCMD(ctx *cli.Context) error {
	return k8s.GetServiceAccount(c.Resources)
}

func (c icli) InitCMD() *cli.Command {
	ef := excludedFlagsFromConfig()
	configFlags := []cli.Flag{}
	for _, vs := range c.flags.all {
		for _, flag := range vs {
			if !slices.Contains(ef, flag.Names()[0]) {
				configFlags = append(configFlags, flag)
			}
		}
	}

	configFlags = append(configFlags, &cli.BoolFlag{
		Name:  persistFlag,
		Value: false,
		Usage: fmt.Sprintf("will write the file content in %s", defaults.ConfigFile),
	})

	return &cli.Command{
		Name:  "init",
		Usage: "create configuration for the cli [EXPERIMENTAL]",
		Subcommands: []*cli.Command{
			{
				Name:   "github",
				Usage:  "create a github pipeline yaml file",
				Flags:  c.CommandFlags([]FlagsType{InitGithub}),
				Action: c.InitGithubCMD,
				Before: c.baseBeforeFunc,
			},
			{
				Name:   "config",
				Usage:  "create a config file with all available flags set to null",
				Flags:  configFlags,
				Action: c.InitConfigCMD,
				Before: func(ctx *cli.Context) error {
					if err := c.loadFlagsFromConfig(ctx); err != nil {
						return err
					}

					ignoredFlags := excludedFlagsFromConfig()

					return c.checkRequiredFlags(ctx, ignoredFlags)
				},
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
