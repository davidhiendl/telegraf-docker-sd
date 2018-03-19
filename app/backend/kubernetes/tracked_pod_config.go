package kubernetes

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (tp *TrackedPod) ConfigGet(key string) string {
	value, ok := tp.Config[key];
	if ok {
		return value
	} else {
		return ""
	}
}

func (tp *TrackedPod) ConfigOrDefault(key string, def string) string {
	value, ok := tp.Config[key];
	if ok {
		return value
	} else {
		return def
	}
}

func (tp *TrackedPod) ConfigExists(key string, def string) bool {
	_, ok := tp.Config[key];
	return ok
}

func (tp *TrackedPod) ConfigEquals(key string, value string) bool {
	return tp.Config[key] == value
}

func (tp *TrackedPod) ConfigMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tp.ConfigOrDefault(key, "")
	return regex.MatchString(val)
}

// extract configuration values from annoations
func (tp *TrackedPod) parseAnnotationsAsConfig() {
	rex, err := regexp.Compile("^" + regexp.QuoteMeta("telegraf.sd.config/") + "(.+)$")
	if err != nil {
		panic(err)
	}

	for key, value := range tp.Pod.Annotations {
		matches := rex.FindAllStringSubmatch(key, -1)
		if matches != nil {
			shortName := matches[0][1]
			tp.Config[shortName] = value
			continue
		}
	}
}
