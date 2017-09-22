package tgtemplate

import (
	"os"
	"strings"
)

type Params struct {
	EnvMap map[string]string
}

// Create new config and populate it from environment
func NewParams() (*Params) {
	params := Params{
		EnvMap: make(map[string]string),
	}

	// convert environment to map
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		params.EnvMap[pair[0]] = pair[1]
	}

	return &params
}
