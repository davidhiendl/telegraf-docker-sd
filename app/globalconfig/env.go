package globalconfig

import (
	"regexp"
	"github.com/sirupsen/logrus"
)

func (gc *GlobalConfigSpec) EnvGet(key string) string {
	if val, ok := gc.EnvMap[key]; ok {
		return val
	}
	return ""
}

func (gc *GlobalConfigSpec) EnvHas(key string) bool {
	_, ok := gc.EnvMap[key];
	return ok
}

func (gc *GlobalConfigSpec) EnvOrDefault(key string, def string) string {
	if val, ok := gc.EnvMap[key]; ok {
		return val
	}
	return def
}

func (gc *GlobalConfigSpec) EnvEquals(key string, value string) bool {
	return gc.EnvMap[key] == value
}

func (gc *GlobalConfigSpec) EnvMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logrus.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := gc.EnvOrDefault(key, "")
	return regex.MatchString(val)
}
