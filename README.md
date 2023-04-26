# k8s Kurated Addons Cli

This project is intended for building an executable file to run on various CI/CD platforms in order to get a project deployed to a containerised environment

### Pre-requisites

You need Docker installed (or similar solutions) to run this project.

Here you can find a list of possible candidates:

- [Docker](https://docs.docker.com/engine/install/) ( cross-platform, paid solution )
- [Rancher Desktop](https://rancherdesktop.io/) ( cross-platform, FOSS )
- [lima](https://github.com/lima-vm/lima) + [nerdctl](https://github.com/containerd/nerdctl) ( macOS only )


This project assumes that you are able to push to your container repository. You can test this by running

```bash

docker push <yourcontainer.repo/imagename>
```

### Updating the pre-built Dockerfile

This project contains a Dockerfile within `docker/Dockerfile`, which acts as the default Dockerfile when none is given. It is built into the compliation of the CLI.

To accomplish this, this project uses go-bindata, which should be installed once you run

``` bash
make install
```

If you would like to update the Dockerfile, make your changes and run the make command:

```bash
make update-dockerfile
```

### Running the executable

In order to build the executable you simply need to run 

```
make
```                                                 |

You will be able to run the executable from 

```bash
./bin/kka-cli
```