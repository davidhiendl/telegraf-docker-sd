package templatedata

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

// extract configuration values from labels
func (td *TemplateData) parseLabelsAsConfig() {
	rex, err := regexp.Compile("^" + regexp.QuoteMeta("telegraf.sd.config.") + "([a-zA-Z0-9_\\.\\-]*)$")
	if err != nil {
		panic(err)
	}

	for label, value := range td.Container.Labels {
		matches := rex.FindAllStringSubmatch(label, -1)
		if matches == nil {
			continue
		}

		shortName := matches[0][1]
		td.Config[shortName] = value
	}
}

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








