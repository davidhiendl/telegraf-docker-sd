package templatedata

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)


func (td *TemplateData) ConfigGet(key string) string {
	value, ok := td.Config[key];
	if ok {
		return value
	} else {
		return ""
	}
}

func (td *TemplateData) ConfigOrDefault(key string, def string) string {
	value, ok := td.Config[key];
	if ok {
		return value
	} else {
		return def
	}
}

func (td *TemplateData) ConfigExists(key string, def string) bool {
	_, ok := td.Config[key];
	return ok
}

func (tds *TemplateData) ConfigEquals(key string, value string) bool {
	return tds.Config[key] == value
}

func (tds *TemplateData) ConfigMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tds.ConfigOrDefault(key, "")
	return regex.MatchString(val)
}







