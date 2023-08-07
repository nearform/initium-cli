package docker

import (
	"fmt"
	"strings"
)

type DockerImage struct {
	Registry  string
	Name      string
	Directory string
	Tag       string
}

func (dockerImage DockerImage) RemoteTag() string {
	return fmt.Sprintf("%s/%s:%s", strings.ToLower(dockerImage.Registry), dockerImage.Name, dockerImage.Tag)
}

func (dockerImage DockerImage) LocalTag() string {
	return fmt.Sprintf("%s:%s", dockerImage.Name, dockerImage.Tag)
}
