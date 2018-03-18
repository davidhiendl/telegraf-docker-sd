package kubernetes

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

type KubernetesConfigSpec struct {
	AutoConfPrefix string `envconfig:"AUTO_CONF_PREFIX",default:"kubernetes_"`
	TagsFromLabels string `envconfig:"TAGS_FROM_LABELS"`
}

func LoadConfig() *KubernetesConfigSpec {
	cfg := &KubernetesConfigSpec{}
	err := envconfig.Process("TSD_KUBERNETES_", cfg)

	if err != nil {
		logger.Fatalf("failed to parse config: %v", err)
	}

	return cfg
}
