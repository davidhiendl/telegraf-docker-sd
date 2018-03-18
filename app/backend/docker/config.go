package docker

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

type DockerConfigSpec struct {
	AutoConfPrefix string `envconfig:"AUTO_CONF_PREFIX",default:"docker_"`
	TagsFromLabels string `envconfig:"TAGS_FROM_LABELS"`
	TagsFromSwarm  bool   `envconfig:"TAGS_FROM_SWARM",default:"true"`
}

func LoadConfig() *DockerConfigSpec {
	cfg := &DockerConfigSpec{}
	err := envconfig.Process("TSD_DOCKER", cfg)

	if err != nil {
		logger.Fatalf("failed to parse config: %v", err)
	}

	return cfg
}
