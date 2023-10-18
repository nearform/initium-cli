package main

import (
	"embed"
	"os"

	"github.com/charmbracelet/log"
	"github.com/nearform/initium-cli/src/cli"
)

//go:embed assets
var resources embed.FS

// Based on https://goreleaser.com/cookbooks/using-main.version/
// The values will be set at release time by the github action
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	icli := cli.New(resources, cli.Release{Version: version, Commit: commit, Date: date})

	if err := icli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
