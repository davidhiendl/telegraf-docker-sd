package sdtemplate

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
	"os"
	"strings"
)

type Params struct {
	Container types.Container
	Bridge    *network.EndpointSettings
	Tags      map[string]string
	Config    map[string]string
	Image     *types.ImageSummary
	EnvMap map[string]string
}

// Create new config and populate it from environment
func NewParams(container types.Container, image *types.ImageSummary) (*Params) {
	params := Params{
		Container: container,
		Bridge:    container.NetworkSettings.Networks["bridge"],
		Tags:      make(map[string]string),
		Config:    make(map[string]string),
		Image:     image,
		EnvMap: make(map[string]string),
	}

	// convert environment to map
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		params.EnvMap[pair[0]] = pair[1]
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
