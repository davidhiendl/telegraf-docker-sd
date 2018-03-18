package docker

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (tds *TrackedContainer) EnvGet(key string) string {
	if val, ok := tds.EnvMap[key]; ok {
		return val
	}
	return ""
}

func (tds *TrackedContainer) EnvHas(key string) bool {
	_, ok := tds.EnvMap[key];
	return ok
}

func (tds *TrackedContainer) EnvOrDefault(key string, def string) string {
	if val, ok := tds.EnvMap[key]; ok {
		return val
	}
	return def
}

func (tds *TrackedContainer) EnvEquals(key string, value string) bool {
	return tds.EnvMap[key] == value
}

func (tds *TrackedContainer) EnvMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tds.EnvOrDefault(key, "")
	return regex.MatchString(val)
}
