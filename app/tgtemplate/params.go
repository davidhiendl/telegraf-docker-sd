package tgtemplate

import (
	"os"
	"strings"
)

type Params struct {
	EnvMap         map[string]string
	DockerLabelMap map[string]bool
}

// Create new config and populate it from environment
func NewParams(dockerLabels map[string]bool) (*Params) {
	params := Params{
		EnvMap:         make(map[string]string),
		DockerLabelMap: dockerLabels,
	}

	// convert environment to map
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		params.EnvMap[pair[0]] = pair[1]
	}

	return &params
}
