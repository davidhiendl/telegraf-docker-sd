package tgtemplate

import "strings"

const GLOBAL_TAG_ENV_PREFIX = "GLOBAL_TAG_"

func (params *Params) GlobalTagFromEnv(key string) string {
	value, ok := params.EnvMap[GLOBAL_TAG_ENV_PREFIX+key];
	if ok {
		return key + " = " + value
	} else {
		return ""
	}
}

func (params *Params) GlobalTagsFromEnv() map[string]string {
	results := make(map[string]string)

	for key, value := range params.EnvMap {
		if (strings.HasPrefix(key, GLOBAL_TAG_ENV_PREFIX)) {
			results[strings.TrimPrefix(key, GLOBAL_TAG_ENV_PREFIX)] = value
		}
	}

	return results
}

func (params *Params) DockerLabelsInclude() []string {
	keys := []string{}
	for k, v := range params.DockerLabelMap {
		if v {
			keys = append(keys, k)
		}
	}

	return keys
}

func toTomlArray(values []string) string {
	if len(values) <= 0 {
		return "[]"
	}

	return `["` + strings.Join(values, `","`) + `"]`
}

func (params *Params) DockerLabelsIncludeAsTomlArray() string {
	return toTomlArray(params.DockerLabelsInclude())
}

func (params *Params) DockerLabelsExclude() []string {
	keys := []string{}
	for k, v := range params.DockerLabelMap {
		if !v {
			keys = append(keys, k)
		}
	}

	return keys
}

func (params *Params) DockerLabelsExcludeAsTomlArray() string {
	return toTomlArray(params.DockerLabelsExclude())
}
