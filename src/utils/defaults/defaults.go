package defaults

import (
	"embed"
)

const ProjectDirectory = "."
const RepoName = "ghcr.io/nearform"
const DockerfileName = "Dockerfile.kka"

const RootDirectoryTests = "../../../"

//go:embed assets
var TemplateFS embed.FS