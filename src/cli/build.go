package cli

import (
	"fmt"
	"path"

	"github.com/nearform/k8s-kurated-addons-cli/src/services/docker"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/logger"
	"github.com/urfave/cli/v2"
)

func (c CLI) Build() *cli.Command {
	return &cli.Command{
		Name:  "build",
		Usage: "build a container image from the project directory",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "runtime-version",
				EnvVars: []string{"KKA_RUNTIME_VERSION"},
			},
		},
		Action: func(cCtx *cli.Context) error {
			repoName := cCtx.String("repo-name")
			dockerFileName := cCtx.String("dockerfile-name")
			project := c.newProject(cCtx)
			docker := docker.New(project, dockerFileName, repoName)

			logger.PrintInfo("Dockerfile Location: " + path.Join(project.Directory, docker.DockerFileName))

			defer project.DeleteDockerFile()
			err := project.AddDockerFile()
			if err != nil {
				return fmt.Errorf("Persisting docker file content: %v", err)
			}

			return docker.Build()
		},
	}
}
