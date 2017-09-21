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
}

// Create new config and populate it from environment
func NewParams(container types.Container) (*Params) {
	params := Params{
		Container: container,
		Bridge:    container.NetworkSettings.Networks["bridge"],
		Tags:      make(map[string]string),
		Config:    make(map[string]string),
	}

	params.parseLabelsAsConfig()
	params.parseExplicitLabelsAsTags()

	return &params
}

func (params *Params) BridgeIP() string {
	return params.Bridge.IPAddress
}
