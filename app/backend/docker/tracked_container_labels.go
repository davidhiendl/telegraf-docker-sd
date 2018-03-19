package docker

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (tc *TrackedContainer) LabelGet(key string) string {
	if val, ok := tc.Container.Labels[key]; ok {
		return val
	}
	return ""
}

func (tc *TrackedContainer) LabelExists(key string) bool {
	_, ok := tc.Container.Labels[key];
	return ok
}

func (tc *TrackedContainer) LabelOrDefault(key string, def string) string {
	if val, ok := tc.Container.Labels[key]; ok {
		return val
	}
	return def
}

func (tc *TrackedContainer) LabelEquals(key string, value string) bool {
	return tc.Container.Labels[key] == value
}

func (tc *TrackedContainer) LabelMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tc.LabelOrDefault(key, "")
	return regex.MatchString(val)
}

func (tc *TrackedContainer) LabelExistsAllOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := tc.Container.Labels[label]; !exists {
			return false
		}
	}

	return true
}

func (tc *TrackedContainer) LabelExistsAnyOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := tc.Container.Labels[label]; exists {
			return true
		}
	}

	return false
}
