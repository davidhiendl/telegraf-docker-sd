package docker

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (tds *TrackedContainer) LabelGet(key string) string {
	if val, ok := tds.Container.Labels[key]; ok {
		return val
	}
	return ""
}

func (tds *TrackedContainer) LabelExists(key string) bool {
	_, ok := tds.Container.Labels[key];
	return ok
}

func (tds *TrackedContainer) LabelOrDefault(key string, def string) string {
	if val, ok := tds.Container.Labels[key]; ok {
		return val
	}
	return def
}

func (tds *TrackedContainer) LabelEquals(key string, value string) bool {
	return tds.Container.Labels[key] == value
}

func (tds *TrackedContainer) LabelMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tds.LabelOrDefault(key, "")
	return regex.MatchString(val)
}

func ( tds *TrackedContainer) LabelExistsAllOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := tds.Container.Labels[label]; !exists {
			return false
		}
	}

	return true
}

func (tds *TrackedContainer) LabelExistsAnyOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := tds.Container.Labels[label]; exists {
			return true
		}
	}

	return false
}
