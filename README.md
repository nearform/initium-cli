# Initium project CLI

A single static binary that can run on any CI to build your code and deploy it in a single step.

All with a nice development workflow in mind like ephemeral environment for your PRs.

### Pre-requisites

1. GoLang

    You can install it with your prefered package manager or usign asfd with `asdf install`.

2. Docker (or similar solutions)  

    Here you can find a list of possible candidates:
    - [Docker](https://docs.docker.com/engine/install/) ( cross-platform, paid solution )
    - [Rancher Desktop](https://rancherdesktop.io/) ( cross-platform, FOSS )
    - [lima](https://github.com/lima-vm/lima) + [nerdctl](https://github.com/containerd/nerdctl) ( macOS only )

### Build the executable

In order to build the executable you simply need to run 

```bash
make build
```

You will be able to run the executable from 

```bash
./bin/initium
```

### Setup local environment

These are the environment variables that you have to set in order to use the onmain, onbranch commands from your local environment

```
export INITIUM_REGISTRY_PASSWORD="<github_pat>"
export INITIUM_REGISTRY_USER="<github_user>"
```

and

```
export INITIUM_CLUSTER_ENDPOINT=$(kubectl config view -o jsonpath='{.clusters[?(@.name == "kind-k8s-kurated-addons")].cluster.server}')
export INITIUM_CLUSTER_TOKEN=$(kubectl get secrets initium-cli-token -o jsonpath="{.data.token}" | base64 -d)
export INITIUM_CLUSTER_CA_CERT=$(kubectl get secrets initium-cli-token -o jsonpath="{.data.ca\.crt}" | base64 -d)
```

### Supported Technologies

Following we have a matrix related to which technologies our cli is currently compatible with. 

| Technologies        | Supported          | Comments    |
|---------------------|:------------------:|:-----------:|
| Node                | :white_check_mark: |             |
| GoLang              | :white_check_mark: |             |
| Python              | :negative_squared_cross_mark:                | Coming Soon |
| Java                | :x:                | Coming Soon |
| C#                  | :x:                | Coming Soon |


### CI Integrations

The matrix below gives an overview of the integration status of our CLI with all major CI platforms.

| CI Systems          | Compatibility      | Comments    |
|---------------------|:------------------:|:-----------:|
| GitHub Actions      | :white_check_mark: |             |
| Gitlab CI           | :x:                | Coming Soon |
| Bitbucket Pipelines | :x:                | Coming Soon |
| CircleCI            | :x:                | Coming Soon |
| Azure Devops        | :x:                | Coming Soon |
| Travis CI           | :x:                | Coming Soon |
| Jenkins             | :x:                | Coming Soon |