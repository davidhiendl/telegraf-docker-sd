package templatedata

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (tds *TemplateData) LabelGet(key string) string {
	if val, ok := tds.Container.Labels[key]; ok {
		return val
	}
	return ""
}

func (tds *TemplateData) LabelExists(key string) bool {
	_, ok := tds.Container.Labels[key];
	return ok
}

func (tds *TemplateData) LabelOrDefault(key string, def string) string {
	if val, ok := tds.Container.Labels[key]; ok {
		return val
	}
	return def
}

func (tds *TemplateData) LabelEquals(key string, value string) bool {
	return tds.Container.Labels[key] == value
}

func (tds *TemplateData) LabelMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tds.LabelOrDefault(key, "")
	return regex.MatchString(val)
}

func ( tds *TemplateData) LabelExistsAllOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := tds.Container.Labels[label]; !exists {
			return false
		}
	}

	return true
}

func (tds *TemplateData) LabelExistsAnyOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := tds.Container.Labels[label]; exists {
			return true
		}
	}

	return false
}
