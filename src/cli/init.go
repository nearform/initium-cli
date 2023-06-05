package cli

import (
	"github.com/nearform/k8s-kurated-addons-cli/src/services/project"
	log "github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"
	"github.com/urfave/cli/v2"
)

func (c CLI) Init(cCtx *cli.Context) error {
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
		log.PrintInfo(v)
	}

	return nil
}

func (c CLI) InitCMD() *cli.Command {
	return &cli.Command{
		Name:  "init",
		Usage: "create a new pipeline file to be used for your provider",
		Subcommands: []*cli.Command{
			{
				Name:   "github",
				Usage:  "create a github pipeline yaml file",
				Flags:  Flags(InitGithub),
				Action: c.Init,
			},
		},
	}
}
