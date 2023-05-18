package docker

import (
    "fmt"
	"github.com/nearform/k8s-kurated-addons-cli/src/utils/defaults"
)

type DockerImage struct {
    Registry    string
    Name        string
    Directory   string
    Tag         string
}

func (di DockerImage) RemoteTag() string {
	if di.Directory != defaults.ProjectDirectory {
		return fmt.Sprintf("%s/%s/%s:%s", di.Registry, di.Name, di.Directory, di.Tag)
	}
	return fmt.Sprintf("%s/%s:%s", di.Registry, di.Name, di.Tag)
}

func (di DockerImage) LocalTag() string {
	if di.Directory != defaults.ProjectDirectory {
		return fmt.Sprintf("%s/%s:%s", di.Name, di.Directory, di.Tag)
	}
	return fmt.Sprintf("%s:%s", di.Name, di.Tag)
}
