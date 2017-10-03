package sdtemplate

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

// execute a custom regex pattern against the container image name
func (params *Params) MatchImageRegex(pattern string) bool {
	expr, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	if expr.MatchString(params.Container.Image) {
		logger.Debugf("matching \"%v\" against \"%v\": true", params.Container.Image, pattern)
		return true
	} else {
		logger.Debugf("matching \"%v\" against \"%v\": false", params.Container.Image, pattern)
	}

	// match against each tag
	for _, tag := range params.Image.RepoTags {
		if expr.MatchString(tag) {
			logger.Debugf("matching \"%v\" against \"%v\": true", tag, pattern)
			return true
		} else {
			logger.Debugf("matching \"%v\" against \"%v\": false", tag, pattern)
		}
	}

	return false
}

// Look for an exact match of the image but ignore the tag
func (params *Params) MatchImage(pattern string) bool {
	return params.MatchImageRegex("^" + regexp.QuoteMeta(pattern) + ":.*$")
}

func (params *Params) LabelExists(label string) bool {
	_, ok := params.Container.Labels[label]
	return ok
}

func (params *Params) LabelExistsAllOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := params.Container.Labels[label]; !exists {
			return false
		}
	}

	return true
}

func (params *Params) LabelExistsAnyOf(labels ...string) bool {
	for _, label := range labels {
		if _, exists := params.Container.Labels[label]; exists {
			return true
		}
	}

	return false
}

func (params *Params) LabelEquals(label string, value string) bool {
	actual, ok := params.Container.Labels[label]
	return ok && actual == value
}
