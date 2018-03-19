package kubernetes

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (tc *TrackedPod) EnvGet(key string) string {
	if val, ok := tc.Env[key]; ok {
		return val
	}
	return ""
}

func (tc *TrackedPod) EnvExists(key string) bool {
	_, ok := tc.Env[key];
	return ok
}

func (tc *TrackedPod) EnvOrDefault(key string, def string) string {
	if val, ok := tc.Env[key]; ok {
		return val
	}
	return def
}

func (tc *TrackedPod) EnvEquals(key string, value string) bool {
	return tc.Env[key] == value
}

func (tc *TrackedPod) EnvMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tc.EnvOrDefault(key, "")
	return regex.MatchString(val)
}
