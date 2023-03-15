# k8s Kurated Addons Cli

This project is intended for building a executable file to run on various CI/CD platfroms in order to get a project deployed to a kubernetes environment

This project assumes that you are able to push to your container repository. You can test this by running

```bash

docker push <yourcontainer.repo/imagename>
```

### Running the executable

In order to build the executable you simply need to run 

```
make
```

The executable takes a few arguments. Running without any arguments will default to values for this project

| Parameter                | Description                                                                                                       |
|--------------------------|-------------------------------------------------------------------------------------------------------------------|
| `--app-name`             | The name of the app. Defaults to `k8s-kurarted-addons-cli `                                                       |
| `--repo-name`            | The base address of the container reppsitory you are wanting to push the image to. Defaults to `ghcr.io/nearform` |
| `--dockerfile-directory` | The directory in which your Dockerfile lives. Defaults to `docker`                                                |       
| `--dockerfile-name`      | The name of the Dockerfile. Defaults to `Dockerfile`                                                              |

You will be able to run the executable from 

```bash
./bin/kka-cli
```