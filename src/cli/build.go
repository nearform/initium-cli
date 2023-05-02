package cli

import (
	"fmt"
	"path"

	"github.com/charmbracelet/log"
	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	"github.com/urfave/cli/v2"
)

func (c CLI) Build(cCtx *cli.Context) error {
	repoName := cCtx.String("repo-name")
	dockerFileName := cCtx.String("dockerfile-name")
	project := c.newProject(cCtx)
	docker, err := docker.New(project, dockerFileName, repoName)
	if err != nil {
		return fmt.Errorf("Creating docker service: %v", err)
	}

	log.Info("Dockerfile Location: " + path.Join(project.Directory, docker.DockerFileName))
	return docker.Build()
}

func (c CLI) BuildCMD() *cli.Command {
	return &cli.Command{
		Name:   "build",
		Usage:  "build a container image from the project directory",
		Flags:  Flags(Build),
		Action: c.Build,
	}
}
