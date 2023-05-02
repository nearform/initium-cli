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
	log.Info("nearForm: k8s kurated addons CLI")

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	cli := cli.CLI{
		Resources: resources,
		CWD:       cwd,
	}

	cli.Run()
}
