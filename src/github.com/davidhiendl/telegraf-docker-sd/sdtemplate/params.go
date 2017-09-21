package sdtemplate

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

type Params struct {
	Container types.Container
	Bridge    *network.EndpointSettings
	Tags      map[string]string
	Config    map[string]string
	Image     *types.ImageSummary
}

// Create new config and populate it from environment
func NewParams(container types.Container, image *types.ImageSummary) (*Params) {
	params := Params{
		Container: container,
		Bridge:    container.NetworkSettings.Networks["bridge"],
		Tags:      make(map[string]string),
		Config:    make(map[string]string),
		Image:     image,
	}

	params.parseLabelsAsConfig()
	params.parseExplicitLabelsAsTags()

	return &params
}

func (params *Params) BridgeIP() string {
	return params.Bridge.IPAddress
}

func (params *Params) Labels() map[string]string {
	return params.Container.Labels
}
