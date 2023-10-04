package cli

import (
	"github.com/nearform/initium-cli/src/services/versions"
	"github.com/nearform/initium-cli/src/utils/logger"
	"github.com/urfave/cli/v2"
)

func (c *icli) Version(cCtx *cli.Context) error {
	cliVersionsFileContent, err := versions.LoadCliVersionsFileContent(c.Resources)
	if err != nil {
		return err
	}

	cliVersion, err := versions.GetCliVersion(cliVersionsFileContent)
	if err != nil {
		return err
	}

	logger.PrintInfo("initium CLI version: " + cliVersion)

	return nil
}

func (c icli) VersionCMD() *cli.Command {
	return &cli.Command{
		Name:   "version",
		Usage:  "show the initium CLI version",
		Flags:  nil,
		Action: c.Version,
		Before: c.baseBeforeFunc,
	}
}
