package cli

import (
	"github.com/docker/docker/api/types"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	"github.com/urfave/cli/v2"
)

func (c CLI) Push(cCtx *cli.Context) error {
	repoName := cCtx.String("repo-name")
	dockerFileName := cCtx.String("dockerfile-name")
	project := c.newProject(cCtx)
	docker := docker.New(project, dockerFileName, repoName)
	docker.AuthConfig = types.AuthConfig{
		Username: cCtx.String("registry-user"),
		Password: cCtx.String("registry-password"),
	}
	return docker.Push()
}

func (c CLI) PushCMD() *cli.Command {
	return &cli.Command{
		Name:   "push",
		Usage:  "push the container image to a registry",
		Flags:  Flags(Registry),
		Action: c.Push,
	}
}
