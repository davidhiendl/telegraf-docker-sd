package config

import (
	"github.com/kelseyhightower/envconfig"
	"path/filepath"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	"strings"
	"os"
	"github.com/davidhiendl/telegraf-docker-sd/app/constants"
)

type ConfigSpec struct {
	TemplateDir       string `envconfig:"TEMPLATE_DIR" default:"/etc/telegraf/sd-tpl.d"`
	ConfigDir         string `envconfig:"CONFIG_DIR" default:"/etc/telegraf/telegraf.d"`
	AutoConfPrefix    string `envconfig:"AUTO_CONF_PREFIX" default:"tsd_"`
	AutoConfExtension string `envconfig:"AUTO_CONF_EXTENSION" default:".conf"`
	LogLevel          string `envconfig:"LOG_LEVEL" default:"warn"`
	CleanOutput       bool   `envconfig:"CLEAN_OUTPUT" default:"true"`
	QueryInterval     int    `envconfig:"QUERY_INTERVAL",default:"15"`

	BackendList string   `envconfig:"BACKENDS"`
	Backends    []string `ignored:"true"`

	GlobalTags map[string]string `ignored:"true"`
	EnvMap     map[string]string `ignored:"true"`
}

func KeyValueMapFromEnv(prefix string, envMap map[string]string) map[string]string {
	results := make(map[string]string)

	for key, value := range envMap {
		if strings.HasPrefix(key, prefix) {
			results[strings.TrimPrefix(key, prefix)] = value
		}
	}

	return results
}

func ConfigListToArray(csv string) []string {
	values := []string{}
	for _, val := range strings.Split(csv, ",") {
		val = strings.Trim(val, " \t\n\r")
		if len(val) > 0 {
			values = append(values, val)
		}
	}

	return values
}

func GetEnvMap() map[string]string {
	// create env map
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		envs[pair[0]] = strings.Trim(pair[1], " \r\n")
	}

	return envs
}

func Load() *ConfigSpec {
	cfg := &ConfigSpec{}
	err := envconfig.Process("TSD", cfg)

	if err != nil {
		logger.Fatalf("failed to parse config: %v", err)
	}

	// set log level
	logger.SetLevel(logger.LevelFromString(cfg.LogLevel))

	// convert csv values
	cfg.Backends = ConfigListToArray(cfg.BackendList)

	// convert paths to absolute
	cfg.TemplateDir, _ = filepath.Abs(cfg.TemplateDir)
	cfg.ConfigDir, _ = filepath.Abs(cfg.ConfigDir)

	cfg.EnvMap = GetEnvMap()
	cfg.GlobalTags = KeyValueMapFromEnv(constants.GLOBAL_TAG_PREFIX, cfg.EnvMap)

	return cfg
}
