package kubernetes

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
)

type KubernetesConfigSpec struct {
	AutoConfPrefix   string `envconfig:"AUTO_CONF_PREFIX",default:"kubernetes_"`
	NodeNameOverride string `envconfig:"NODE_NAME_OVERRIDE"`

	TagNamespace bool `envconfig:"TAG_NAMESPACE" default:"true"`
	TagPodName   bool `envconfig:"TAG_POD_NAME" default:"true"`

	TagLabelsWhitelistStr string `envconfig:"TAG_LABELS_WHITELIST"`
	TagLabelsBlacklistStr string `envconfig:"TAG_LABELS_BLACKLIST"`

	TagLabelsWhitelist []string `ignored:"true"`
	TagLabelsBlacklist []string `ignored:"true"`
}

func LoadConfig() *KubernetesConfigSpec {
	cfg := &KubernetesConfigSpec{}
	err := envconfig.Process("TSD_KUBERNETES", cfg)

	if err != nil {
		logger.Fatalf(LOG_PREFIX+" failed to parse config: %v", err)
	}

	// convert list to array
	cfg.TagLabelsWhitelist = config.ConfigListToArray(cfg.TagLabelsWhitelistStr)
	cfg.TagLabelsBlacklist = config.ConfigListToArray(cfg.TagLabelsBlacklistStr)

	if len(cfg.TagLabelsWhitelist) > 0 && len(cfg.TagLabelsBlacklist) > 0 {
		logger.Fatalf(LOG_PREFIX+" cannot have label whitelist and blacklist", err)
	}

	return cfg
}
