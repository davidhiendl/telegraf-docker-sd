package kubernetes

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (tp *TrackedPod) AnnotationGet(key string) string {
	value, ok := tp.Pod.Annotations[key];
	if ok {
		return value
	} else {
		return ""
	}
}

func (tp *TrackedPod) AnnotationOrDefault(key string, def string) string {
	value, ok := tp.Pod.Annotations[key];
	if ok {
		return value
	} else {
		return def
	}
}

func (tp *TrackedPod) AnnotationExists(key string, def string) bool {
	_, ok := tp.Pod.Annotations[key];
	return ok
}

func (tp *TrackedPod) AnnotationEquals(key string, value string) bool {
	return tp.Pod.Annotations[key] == value
}

func (tp *TrackedPod) AnnotationMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logger.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tp.AnnotationOrDefault(key, "")
	return regex.MatchString(val)
}
