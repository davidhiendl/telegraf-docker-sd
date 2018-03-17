package templatedata

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (tds *TemplateData) EnvGet(key string) string {
	if val, ok := tds.EnvMap[key]; ok {
		return val
	}
	return ""
}

func (tds *TemplateData) EnvHas(key string) bool {
	_, ok := tds.EnvMap[key];
	return ok
}

func (tds *TemplateData) EnvOrDefault(key string, def string) string {
	if val, ok := tds.EnvMap[key]; ok {
		return val
	}
	return def
}

func (tds *TemplateData) EnvEquals(key string, value string) bool {
	return tds.EnvMap[key] == value
}

func (tds *TemplateData) EnvMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tds.EnvOrDefault(key, "")
	return regex.MatchString(val)
}
