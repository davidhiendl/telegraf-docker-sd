package docker

import (
	"github.com/docker/docker/api/types"
	td "github.com/davidhiendl/telegraf-docker-sd/app/backend/docker/templatedata"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	"path/filepath"
)

// TrackedContainer is used to maintain state about already processed containers and to be able to remove their configurations easily
type TrackedContainer struct {
	ID         string
	ShortID    string
	backend    *DockerBackend
	container  *types.Container
	configFile string
	Data       *td.TemplateData
}

// Create new config and populate it from environment
func NewTrackedContainer(backend *DockerBackend, container *types.Container) *TrackedContainer {
	tc := TrackedContainer{
		ID:        container.ID,
		ShortID:   toShortID(container.ID),
		backend:   backend,
		container: container,
		Data:      td.NewTemplateData(container),
	}

	tc.Data.Image = tc.backend.getImageForID(container.ImageID)

	// add explicit labels
	tc.parseExplicitLabelsAsTags()

	// add swarm labels if desired
	if tc.backend.config.TagsFromSwarm {
		tc.parseLabelsAsTags(SWARM_LABELS)
	}

	// debug
	logger.Debugf("[docker][%v] tags: %+v", tc.ShortID, tc.Data.Tags)
	logger.Debugf("[docker][%v] config: %+v", tc.ShortID, tc.Data.Config)
	logger.Debugf("[docker][%v] labels: %+v", tc.ShortID, tc.Data.Container.Labels)

	return &tc
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
