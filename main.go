package main

import (
	"embed"

	"github.com/charmbracelet/log"
	"github.com/nearform/k8s-kurated-addons-cli/src/cli"
)

//go:embed assets
var resources embed.FS

func main() {
	log.Info("nearForm: k8s kurated addons CLI")

	cli := cli.CLI{
		Resources: resources
	}

	cli.Run()
}
