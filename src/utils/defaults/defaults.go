package defaults

const (
	ProjectDirectory    string = "."
	RepoName            string = "ghcr.io/nearform"
	GithubActionFolder  string = ".github/workflows"
	GithubDefaultBranch string = "main"
	ConfigFile          string = ".initium.yaml"
	AppVersion          string = "latest"
	GeneratedDockerFile string = "Dockerfile.initium"
	EnvVarFile          string = ".env.initium"
)

// renovate: datasource=docker depName=node
const DefaultNodeRuntimeVersion string = "20.2.0"

// renovate: datasource=docker depName=golang
const DefaultGoRuntimeVersion string = "1.20.4"

const DefaultFrontendJsRuntimeVersion string = "20.2.0"
