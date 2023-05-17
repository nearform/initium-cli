package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *CLI) template(cCtx *cli.Context) error {
	project := c.newProject(cCtx)
	content, err := project.Dockerfile()
	if err != nil {
		return fmt.Errorf("Getting docker file %v", err)
	}
	fmt.Println(string(content))
	return nil
}

func (c *CLI) TemplateCMD() *cli.Command {
	return &cli.Command{
		Name:   "template",
		Usage:  "output the docker file used for this project",
		Flags:  Flags(Build),
		Action: c.template,
	}
}
