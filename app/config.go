package app

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	"path/filepath"
	"strconv"
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
	CleanOutput         bool   `envconfig:"CLEAN_OUTPUT"`
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
	CONFIG_DEFAULT_QUERY_INTERVAL = 15
)

var DefaultConfig = Config{
	TemplateDir:         "/etc/telegraf/sd-tpl.d",
	ConfigDir:           "/etc/telegraf/telegraf.d",
	AutoConfPrefix:      "sd-container_",
	AutoConfExtension:   ".conf",
	LogLevel:            logger.LOG_INFO,
	TagsFromSwarmLabels: true,
	TagsFromLabels:      "",
	QueryInterval:       CONFIG_DEFAULT_QUERY_INTERVAL,
	CleanOutput:         false,
}

// Create new config and populate it from environment
func NewConfig() (*Config) {
	c := DefaultConfig
	return &c
}

func (c *Config) AsMap() (map[string]string) {
	m := make(map[string]string)

	m["TemplateDir"] = c.TemplateDir
	m["TemplateDir"] = c.ConfigDir
	m["AutoConfPrefix"] = c.AutoConfPrefix
	m["AutoConfExtension"] = c.AutoConfExtension
	m["LogLevel"] = strconv.Itoa(c.LogLevel)
	m["TagsFromSwarmLabels"] = strconv.FormatBool(c.TagsFromSwarmLabels)
	m["TagsFromLabels"] = c.TagsFromLabels
	m["QueryInterval"] = strconv.Itoa(c.QueryInterval)
	m["CleanOutput"] = strconv.FormatBool(c.CleanOutput)

	return m
}

func (c *Config) LoadFromEnv() error {
	err := envconfig.Process("TSD", c)
	c.TemplateDir, _ = filepath.Abs(c.TemplateDir)
	c.ConfigDir, _ = filepath.Abs(c.ConfigDir)
	return err;
}
