package sdtemplate

import (
	"regexp"
	"fmt"
)

// execute a custom regex pattern against the container image name
func (params *Params) MatchImageRegex(pattern string) bool {
	expr, err := regexp.Compile(pattern)
	fmt.Printf("pattern: %v image: %v", pattern, params.Image.RepoTags)
	if err != nil {
		panic(err)
	}

	// match against each tag
	for _, tag := range params.Image.RepoTags {
		if expr.MatchString(tag) {
			fmt.Printf(" matched \n")
			return true
		}
	}
	fmt.Printf(" no match \n")

	return false
}

// Look for an exact match of the image but ignore the tag
func (params *Params) MatchImage(pattern string) bool {
	return params.MatchImageRegex("^" + regexp.QuoteMeta(pattern) + ":.*$")
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
