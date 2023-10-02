package cli

import (
	"fmt"

	"github.com/urfave/cli/v2"
)

func (c *icli) template(cCtx *cli.Context) error {
	project, err := c.getProject(cCtx)
	if err != nil {
		return err
	}
	content, err := project.Dockerfile()
	if err != nil {
		return fmt.Errorf("Getting docker file %v", err)
	}
	fmt.Fprintln(c.Writer, string(content))
	return nil
}

func (c *icli) TemplateCMD() *cli.Command {
	return &cli.Command{
		Name:   "template",
		Usage:  "output the docker file used for this project",
		Flags:  c.CommandFlags([]FlagsType{Build}),
		Action: c.template,
		Before: c.baseBeforeFunc,
	}
}
