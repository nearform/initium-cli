package main

import (
	"embed"

	"github.com/nearform/k8s-kurated-addons-cli/src/cli"
)

//go:embed assets
var resources embed.FS

func main() {
	cli := cli.CLI{
		Resources: resources,
	}
	cli.Run()
}
