package main

import (
	"embed"
	"os"

	"github.com/charmbracelet/log"
	"github.com/nearform/initium-cli/src/cli"
)

//go:embed assets
var resources embed.FS

func main() {
	icli := cli.New(resources)

	if err := icli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
