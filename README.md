# k8s Kurated Addons Cli

This project is intended for building an executable file to run on various CI/CD platforms in order to get a project deployed to a containerised environment

### Pre-requisites

You need Docker installed (or similar solutions) to run this project.

Here you can find a list of possible candidates:

- [Docker](https://docs.docker.com/engine/install/) ( cross-platform, paid solution )
- [Rancher Desktop](https://rancherdesktop.io/) ( cross-platform, FOSS )
- [lima](https://github.com/lima-vm/lima) + [nerdctl](https://github.com/containerd/nerdctl) ( macOS only )


### Running the executable

In order to build the executable you simply need to run 

```
make
```                                                 |

You will be able to run the executable from 

```bash
./bin/kka-cli
```

### Testing

1. Move into the example folder and either run `../bin/kka-cli build` or `go run ../main.go build`
2. Run `docker run ghcr.io/nearform/k8s-kurated-addons-cli:latest`
3. You should see the `KKA-CLI from NearForm` output in your console
4. Remove the image `docker image rmi -f ghcr.io/nearform/k8s-kurated-addons-cli:latest`