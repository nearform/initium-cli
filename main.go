package main

import (
	"embed"
	"os"

	"github.com/charmbracelet/log"
	"github.com/nearform/k8s-kurated-addons-cli/src/cli"
)

//go:embed assets
var resources embed.FS

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cli := cli.CLI{
		Resources: resources,
		CWD:       cwd,
		Logger: log.NewWithOptions(os.Stderr, log.Options{
			Level:           log.ParseLevel(os.Getenv("KKA_LOG_LEVEL")),
			ReportCaller:    true,
			ReportTimestamp: true,
		}),
	}

	cli.Run()
}
