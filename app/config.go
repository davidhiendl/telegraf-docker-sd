package app

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	TemplateDir         string `envconfig:"TEMPLATE_DIR"`
	ConfigDir           string `envconfig:"CONFIG_DIR"`
	AutoConfPrefix      string `envconfig:"AUTO_CONF_PREFIX"`
	AutoConfExtension   string `envconfig:"AUTO_CONF_EXTENSION"`
	LogLevel            int    `envconfig:"LOG_LEVEL"` // TODO fetch this properly from env and also allow non-int constants to be used (NONE;ERROR;WARN;INFO;DEBUG)
	TagsFromSwarmLabels bool   `envconfig:"TAG_SWARM_LABELS"`
	TagsFromLabels      string `envconfig:"TAG_LABELS"`
	QueryInterval       int    `envconfig:"QUERY_INTERVAL"`
}

var SWARM_LABELS = []string{
	"com.docker.stack.namespace",
	"com.docker.swarm.node.id",
	"com.docker.swarm.service.id",
	"com.docker.swarm.service.name",
	"com.docker.swarm.task",
	"com.docker.swarm.task.id",
	"com.docker.swarm.task.name",
}

const (
	LOG_NONE  = 1 << iota
	LOG_ERROR = 1 << iota
	LOG_WARN  = 1 << iota
	LOG_INFO  = 1 << iota
	LOG_DEBUG = 1 << iota
)

const (
	CONFIG_DEFAULT_QUERY_INTERVAL = 15
)

var DefaultConfig = Config{
	TemplateDir:         "/etc/telegraf/sd-tpl.d",
	ConfigDir:           "/etc/telegraf/telegraf.d",
	AutoConfPrefix:      "sd-container_",
	AutoConfExtension:   ".conf",
	LogLevel:            LOG_WARN,
	TagsFromSwarmLabels: true,
	TagsFromLabels:      "",
	QueryInterval:       CONFIG_DEFAULT_QUERY_INTERVAL,
}

// Create new config and populate it from environment
func NewConfig() (*Config) {
	c := DefaultConfig
	return &c
}

func (c *Config) LoadFromEnv() error {
	err := envconfig.Process("TSD", c)
	return err;
}
