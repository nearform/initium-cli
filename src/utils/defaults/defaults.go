package defaults

const ProjectDirectory = "."
const RepoName = "ghcr.io/nearform"
const DockerfileName = "Dockerfile.kka"
const GithubActionFolder = ".github/workflows"
const GithubDefaultBranch = "main"
const ConfigFile = ".kka"

// renovate: datasource=docker depName=node
const DefaultNodeRuntimeVersion = "20.2.0"

// renovate: datasource=docker depName=golang
const DefaultGoRuntimeVersion = "1.20.4"
