package docker

import (
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

// execute a custom regex pattern against the container image name
func (tc *TrackedContainer) MatchImageRegex(pattern string) bool {
	expr, err := regexp.Compile(pattern)
	if err != nil {
		panic(err)
	}

	if expr.MatchString(tc.Container.Image) {
		logger.Debugf(`matching image = "%v" against pattern = "%v": true`, tc.Container.Image, pattern)
		return true
	} else {
		logger.Debugf(`matching image = "%v" against pattern = "%v": false`, tc.Container.Image, pattern)
	}

	// match against each tag
	for _, tag := range tc.Image.RepoTags {
		if expr.MatchString(tag) {
			logger.Debugf(`matching image = "%v" against pattern = "%v": true`, tag, pattern)
			return true
		} else {
			logger.Debugf(`matching image = "%v" against pattern = "%v": false`, tag, pattern)
		}
	}

	return false
}

// Look for an exact match of the image but ignore the tag
func (tc *TrackedContainer) MatchImage(pattern string) bool {
	return tc.MatchImageRegex("^" + regexp.QuoteMeta(pattern) + "(:.*)?$")
}
