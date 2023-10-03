package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/charmbracelet/log"
	"k8s.io/utils/strings/slices"

	"github.com/nearform/initium-cli/src/services/k8s"
	"github.com/nearform/initium-cli/src/services/project"
	"github.com/nearform/initium-cli/src/utils/defaults"
	"github.com/urfave/cli/v2"

	knative "github.com/nearform/initium-cli/src/services/k8s"
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

func (c icli) InitConfigCMD(ctx *cli.Context) error {
	excludedFlags := []string{
		"help",
		appVersionFlag,
		namespaceFlag,
		configFileFlag,
		projectDirectoryFlag,
		destinationFolderFlag,
		tokenFlag,
		registryPasswordFlag,
		caCRTFlag,
	}

	f := []cli.Flag{}
	for _, vs := range c.flags.all {
		f = append(f, vs...)
	}

	sort.Sort(cli.FlagsByName(f))

	config := ""
	for _, flag := range f {
		stringFlag := flag.(*cli.StringFlag)
		if slices.Contains(excludedFlags, stringFlag.Name) {
			continue
		}

		value := ctx.String(stringFlag.Name)
		if value == "" {
			value = stringFlag.Value
		}

		next := ""
		if value == "" {
			next = fmt.Sprintf("%s: null\n", stringFlag.Name)
		} else {
			next = fmt.Sprintf("%s: %s\n", stringFlag.Name, value)
		}

		config = config + next
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

func (c icli) InitKnativeDomainCMD(cCtx *cli.Context) error {
	config, err := knative.Config(
		cCtx.String(endpointFlag),
		cCtx.String(tokenFlag),
		[]byte(cCtx.String(caCRTFlag)),
	)

	if err != nil {
		return err
	}

	if cCtx.NArg() < 1 {
		return fmt.Errorf("Knative domain argument is required!")
	}

	knativeDomain := cCtx.Args().Get(0)

	err = knative.DomainUpd(knativeDomain, config)

	if err != nil {
		return err
	}

	return nil
}

func (c icli) InitCMD() *cli.Command {
	configFlags := c.CommandFlags([]FlagsType{Shared})
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
				Before: c.baseBeforeFunc,
			},
			{
				Name:   "service-account",
				Usage:  "output all resources needed to create a service account",
				Action: c.InitServiceAccountCMD,
				Before: c.baseBeforeFunc,
			},
			{
				Name:   "knative-domain",
				Usage:  "updates knative service default domain",
				Flags:  c.CommandFlags([]FlagsType{Kubernetes}),
				Action: c.InitKnativeDomainCMD,
				Before: func(ctx *cli.Context) error {
					if err := c.loadFlagsFromConfig(ctx); err != nil {
						return err
					}
		
					ignoredFlags := []string{namespaceFlag}
		
					return c.checkRequiredFlags(ctx, ignoredFlags)
				},
			},
		},
	}
}
