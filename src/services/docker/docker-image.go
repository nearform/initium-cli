package docker

import (
	"fmt"
)

type DockerImage struct {
	Registry  string
	Name      string
	Directory string
	Tag       string
}

func (dockerImage DockerImage) RemoteTag() string {
	return fmt.Sprintf("%s/%s:%s", dockerImage.Registry, dockerImage.Name, dockerImage.Tag)
}

func (dockerImage DockerImage) LocalTag() string {
	return fmt.Sprintf("%s:%s", dockerImage.Name, dockerImage.Tag)
}
