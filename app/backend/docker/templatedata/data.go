package templatedata

import (
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/network"
	"strings"
)

type TemplateData struct {
	Name      string
	Container *types.Container
	EnvMap    map[string]string
	Tags      map[string]string
	Config    map[string]string
	Image     *types.ImageSummary
}

func NewTemplateData(container *types.Container) *TemplateData {
	td := &TemplateData{
		Tags:      make(map[string]string),
		Container: container,
		Name:      strings.TrimLeft(container.Names[0], "/"),
	}

	td.parseLabelsAsConfig()

	return td
}

func (td *TemplateData) ShortID() string {
	return td.Container.ID[0:12]
}

func (td *TemplateData) dockerNetBridge() *network.EndpointSettings {
	return td.Container.NetworkSettings.Networks["bridge"]
}

func (td *TemplateData) BridgeIP() string {
	return td.dockerNetBridge().IPAddress
}
