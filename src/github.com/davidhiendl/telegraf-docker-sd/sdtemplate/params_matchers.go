package sdtemplate

import (
	"regexp"
)

// execute a custom regex pattern against the container image name
func (params *Params) MatchImageRegex(pattern string) bool {
	match, err := regexp.MatchString(pattern, params.Container.Image)

	if err != nil {
		panic(err)
	}

	return match
}

// Look for an exact match of the image but ignore the tag
func (params *Params) MatchImage(pattern string) bool {
	match, err := regexp.MatchString("^"+regexp.QuoteMeta(pattern)+":.*$", params.Container.Image)

	if err != nil {
		panic(err)
	}

	return match
}

func (params *Params) LabelExists(label string) bool {
	_, ok := params.Container.Labels[label];
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
	actual, ok := params.Container.Labels[label];
	return ok && actual == value
}
