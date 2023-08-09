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


### Run from code

1. Run `go run main.go --project-directory example build`
2. Run `docker run ghcr.io/nearform/initium-cli:latest`
3. You should see the `KKA-CLI from NearForm` output in your console
4. Remove the image `docker image rmi -f ghcr.io/nearform/k8s-kurated-addons-cli:latest`

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
export KKA_REGISTRY_PASSWORD="<github_pat>"
export KKA_REGISTRY_USER="<github_user>"
```

and

```
export KKA_CLUSTER_ENDPOINT=$(kubectl config view -o jsonpath='{.clusters[?(@.name == "kind-k8s-kurated-addons")].cluster.server}')
export KKA_CLUSTER_TOKEN=$(kubectl get secrets initium-cli-token -o jsonpath="{.data.token}" | base64 -d)
export KKA_CLUSTER_CA_CERT=$(kubectl get secrets initium-cli-token -o jsonpath="{.data.ca\.crt}" | base64 -d)
```
