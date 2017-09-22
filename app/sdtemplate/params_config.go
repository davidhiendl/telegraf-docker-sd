package sdtemplate

import (
	"regexp"
)

// extract configuration values from labels
func (params *Params) parseLabelsAsConfig() {
	rex, err := regexp.Compile("^telegraf\\.sd\\.config\\.([a-zA-Z0-9_\\.\\-]*)$")
	if err != nil {
		panic(err)
	}

	for label, value := range params.Container.Labels {
		matches := rex.FindAllStringSubmatch(label, -1)
		if matches == nil {
			continue
		}

		shortName := matches[0][1]
		params.Config[shortName] = value
	}
}

func (params *Params) ConfigOrDefault(key string, def string) string {
	value, ok := params.Config[key];
	if ok {
		return value
	} else {
		return def
	}
}

func (params *Params) ConfigGet(key string) string {
	value, ok := params.Config[key];
	if ok {
		return value
	} else {
		return ""
	}
}

func (params *Params) ConfigHas(key string, def string) bool {
	_, ok := params.Config[key];
	return ok
}
