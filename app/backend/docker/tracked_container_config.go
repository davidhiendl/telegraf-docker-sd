package docker

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

// extract configuration values from labels
func (tc *TrackedContainer) parseLabelsAsConfig() {
	rex, err := regexp.Compile("^" + regexp.QuoteMeta("telegraf.sd.config.") + "(.+)$")
	if err != nil {
		panic(err)
	}

	for label, value := range tc.Container.Labels {
		matches := rex.FindAllStringSubmatch(label, -1)
		if matches == nil {
			continue
		}

		shortName := matches[0][1]
		tc.Config[shortName] = value
	}
}

func (tc *TrackedContainer) ConfigGet(key string) string {
	value, ok := tc.Config[key];
	if ok {
		return value
	} else {
		return ""
	}
}

func (tc *TrackedContainer) ConfigOrDefault(key string, def string) string {
	value, ok := tc.Config[key];
	if ok {
		return value
	} else {
		return def
	}
}

func (tc *TrackedContainer) ConfigExists(key string, def string) bool {
	_, ok := tc.Config[key];
	return ok
}

func (tc *TrackedContainer) ConfigEquals(key string, value string) bool {
	return tc.Config[key] == value
}

func (tc *TrackedContainer) ConfigMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tc.ConfigOrDefault(key, "")
	return regex.MatchString(val)
}
