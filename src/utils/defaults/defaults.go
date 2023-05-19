package defaults

const ProjectDirectory = "."
const RepoName = "ghcr.io/nearform"
const DockerfileName = "Dockerfile.kka"

//  renovate: datasource=docker depName=node versioning=node
const DefaultNodeRuntimeVersion = "20"

//  renovate: datasource=docker depName=golang versioning=gomod
const DefaultGoRuntimeVersion = "1.19.3"
