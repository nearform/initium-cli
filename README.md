# k8s Kurated Addons Cli

This project is intended for building a executable file to run on various CI/CD platfroms in order to get a project deployed to a kubernetes environment

### Running the executable

In order to run the executable you simply need to run 

```
make
```

#### Interactive

```bash
./bin/kka-cli
```

This will launch an interactive prompt in which you will need to provide a set of parameters to get your integration working correctly 

#### Non-Interactive

To use the non-interactive (automated) version of this tool, you will need to specify the flag `--non-interactive` or `-ni`. You will need to provide the application with a series of environment variables or parameters

`KKA_APP_NAME` or `--app-name`: The application name

`KKA_APP_PORT` or `--app-port`: The port in which the application runs on
