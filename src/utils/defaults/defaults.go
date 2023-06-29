package defaults

const (
	ProjectDirectory    string = "."
	RepoName            string = "ghcr.io/nearform"
	DockerfileName      string = "Dockerfile.kka"
	GithubActionFolder  string = ".github/workflows"
	GithubDefaultBranch string = "main"
	ConfigFile          string = ".kka"
	AppVersion          string = "latest"
)

// renovate: datasource=docker depName=node
const DefaultNodeRuntimeVersion string = "20.2.0"

// renovate: datasource=docker depName=golang
const DefaultGoRuntimeVersion string = "1.20.4"
