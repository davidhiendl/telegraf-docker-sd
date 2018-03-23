package docker

import (
	"github.com/docker/docker/api/types"
	"path/filepath"
	"strings"
	"github.com/docker/docker/api/types/network"
	"github.com/sirupsen/logrus"
)

// TrackedContainer is used to maintain state about already processed containers and to be able to remove their configurations easily
type TrackedContainer struct {
	ID      string
	ShortID string
	Name    string

	backend    *DockerBackend
	configFile string

	Container *types.Container
	Env       map[string]string
	Tags      map[string]string
	Config    map[string]string
	Image     *types.ImageSummary
}

// Create new config and populate it from environment
func NewTrackedContainer(backend *DockerBackend, container *types.Container) *TrackedContainer {
	tc := TrackedContainer{
		backend:   backend,
		Container: container,

		ID:      container.ID,
		ShortID: toShortID(container.ID),
		Name:    strings.TrimLeft(container.Names[0], "/"),

		Config: make(map[string]string),
		Tags:   make(map[string]string),
		Env: backend.commonConfig.EnvMap,
	}

	tc.Image = tc.backend.getImageForID(container.ImageID)

	// add explicit labels
	tc.parseLabelsAsTags()

	// add swarm labels if desired
	if tc.backend.config.TagsFromSwarm {
		tc.parseSwarmLabelsAsTags()
	}

	// debug
	logrus.Debugf(LOG_PREFIX+"[%v] tags: %+v", tc.ShortID, tc.Tags)
	logrus.Debugf(LOG_PREFIX+"[%v] config: %+v", tc.ShortID, tc.Config)
	logrus.Debugf(LOG_PREFIX+"[%v] labels: %+v", tc.ShortID, tc.Container.Labels)

	return &tc
}

func (tc *TrackedContainer) dockerNetBridge() *network.EndpointSettings {
	return tc.Container.NetworkSettings.Networks["bridge"]
}

func (tc *TrackedContainer) BridgeIP() string {
	return tc.dockerNetBridge().IPAddress
}

func (tc *TrackedContainer) GetConfigFile() string {
	if tc.configFile == "" {
		file, _ := filepath.Abs(
			tc.backend.commonConfig.ConfigDir +
				"/" + tc.backend.commonConfig.AutoConfPrefix +
				tc.backend.Name() + "_" +
				tc.ShortID +
				tc.backend.commonConfig.AutoConfExtension)
		tc.configFile = file
	}
	return tc.configFile
}
