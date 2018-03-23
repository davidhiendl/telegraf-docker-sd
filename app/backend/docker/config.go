package docker

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
	"github.com/sirupsen/logrus"
)

type DockerConfigSpec struct {
	AutoConfPrefix string `envconfig:"AUTO_CONF_PREFIX",default:"docker_"`
	TagsFromSwarm  bool   `envconfig:"TAGS_FROM_SWARM",default:"true"`

	TagLabelsWhitelistStr string `envconfig:"TAG_LABELS_WHITELIST"`
	TagLabelsBlacklistStr string `envconfig:"TAG_LABELS_BLACKLIST"`

	TagLabelsWhitelist []string `ignored:"true"`
	TagLabelsBlacklist []string `ignored:"true"`
}

func LoadConfig() *DockerConfigSpec {
	cfg := &DockerConfigSpec{}
	err := envconfig.Process("TSD_DOCKER", cfg)

	if err != nil {
		logrus.Fatalf("failed to parse config: %v", err)
	}

	// convert list to array
	cfg.TagLabelsWhitelist = config.ConfigListToArray(cfg.TagLabelsWhitelistStr)
	cfg.TagLabelsBlacklist = config.ConfigListToArray(cfg.TagLabelsBlacklistStr)

	if len(cfg.TagLabelsWhitelist) > 0 && len(cfg.TagLabelsBlacklist) > 0 {
		logrus.Fatalf(LOG_PREFIX+" cannot have label whitelist and blacklist", err)
	}

	return cfg
}
