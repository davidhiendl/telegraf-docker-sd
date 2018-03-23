package kubernetes

import (
	"regexp"
	"github.com/sirupsen/logrus"
)

func (tc *TrackedPod) LabelGet(key string) string {
	if val, ok := tc.Pod.Labels[key]; ok {
		return val
	}
	return ""
}

func (tc *TrackedPod) LabelExists(key string) bool {
	_, ok := tc.Pod.Labels[key];
	return ok
}

func (tc *TrackedPod) LabelOrDefault(key string, def string) string {
	if val, ok := tc.Pod.Labels[key]; ok {
		return val
	}
	return def
}

func (tc *TrackedPod) LabelEquals(key string, value string) bool {
	return tc.Pod.Labels[key] == value
}

func (tc *TrackedPod) LabelMatches(key string, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		logrus.Errorf("failed to compile template regex: %v" + pattern)
	}

	val := tc.LabelOrDefault(key, "")
	return regex.MatchString(val)
}

func (tc *TrackedPod) LabelExistsAllOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := tc.Pod.Labels[label]; !exists {
			return false
		}
	}

	return true
}

func (tc *TrackedPod) LabelExistsAnyOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := tc.Pod.Labels[label]; exists {
			return true
		}
	}

	return false
}
