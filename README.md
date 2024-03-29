# Initium project CLI

A single static binary that can run on any CI to build your code and deploy it in a single step.

All with a nice development workflow in mind like ephemeral environment for your PRs.

### Running the CLI

The CLI is currently published on [NPM](https://www.npmjs.com/package/initium). You can run it using

```
npx initium
```

That will install the Initium CLI in your folder's `node_modules`.

We recommend you not to try installing the CLI globally if you don't use a Node.js version manager (like `nvm`), as you can run into permission errors.

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
export INITIUM_CLUSTER_ENDPOINT=$(kubectl config view -o jsonpath='{.clusters[?(@.name == "kind-initium-platform")].cluster.server}')
export INITIUM_CLUSTER_TOKEN=$(kubectl get secrets initium-cli-token -o jsonpath="{.data.token}" | base64 -d)
export INITIUM_CLUSTER_CA_CERT=$(kubectl get secrets initium-cli-token -o jsonpath="{.data.ca\.crt}" | base64 -d)
```

### Supported Application Runtimes

Following we have a matrix related to which application runtime our CLI is currently compatible with. For each one of them a Dockerfile template is being used in order to provide an easy way to build and deploy the application to a Kubernetes cluster. 

| Application Runtimes | Supported          |
|----------------------|:------------------:|
| [Nodejs](https://github.com/nearform/initium-cli/blob/main/assets/docker/Dockerfile.node.tmpl) | :white_check_mark: |
| [GoLang](https://github.com/nearform/initium-cli/blob/main/assets/docker/Dockerfile.go.tmpl) | :white_check_mark: |
| Python               | Coming Soon        |
| More will be added...|                    |


### CI Integrations

The matrix below gives an overview of the integration status of our CLI with CI platforms. CLI is able to create automatically `build` and `deploy` pipelines in order to enable CICD process for the application it is used with. Related template files are being used to cover different process steps. 

| CI Systems           | Supported          |
|----------------------|:------------------:|
| [GitHub Actions](https://github.com/nearform/initium-cli/tree/main/assets/github) | :white_check_mark: |
| Gitlab CI            | Coming Soon        |
| Azure Devops         | Coming Soon        |
| More will be added...|                    |


Encrypt your secret with:

```
./initium secrets encrypt --publicKey age1zmh77nlvddsz55q5l67d4ufwewvyhentlku9z90t969szd2lnghslnlese --secret <your secret>
```

or if it's not a string

```
./initium secrets encrypt --publicKey age1zmh77nlvddsz55q5l67d4ufwewvyhentlku9z90t969szd2lnghslnlese --base64secret <your secret in base64>
```
