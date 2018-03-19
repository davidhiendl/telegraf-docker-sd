package docker

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (tc *TrackedContainer) EnvGet(key string) string {
	if val, ok := tc.Env[key]; ok {
		return val
	}
	return ""
}

func (tc *TrackedContainer) EnvExists(key string) bool {
	_, ok := tc.Env[key];
	return ok
}

func (tc *TrackedContainer) EnvOrDefault(key string, def string) string {
	if val, ok := tc.Env[key]; ok {
		return val
	}
	return def
}

func (tc *TrackedContainer) EnvEquals(key string, value string) bool {
	return tc.Env[key] == value
}

func (tc *TrackedContainer) EnvMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tc.EnvOrDefault(key, "")
	return regex.MatchString(val)
}
