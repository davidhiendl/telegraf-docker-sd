package docker

import "regexp"

// extract configuration values from labels
func (tc *TrackedContainer) parseExplicitLabelsAsTags() {
	rex, err := regexp.Compile("^" + regexp.QuoteMeta("telegraf.sd.tags.") + "([a-zA-Z0-9_\\.\\-]*)$")
	if err != nil {
		panic(err)
	}

	for label, value := range tc.container.Labels {
		matches := rex.FindAllStringSubmatch(label, -1)
		if matches != nil {
			shortName := matches[0][1]
			tc.Data.Tags[shortName] = value
			continue
		}
	}
}

func (tc *TrackedContainer) parseLabelsAsTags(labelsAsTags []string) {
	for label, value := range tc.container.Labels {
		for _, matchLabel := range labelsAsTags {
			if label == matchLabel {
				tc.Data.Tags[label] = value
				break
			}
		}
	}
}

/*
func (params *Params) DockerLabelsInclude() []string {
	keys := []string{}
	for k, v := range params.DockerLabelMap {
		if v {
			keys = append(keys, k)
		}
	}

	return keys
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
*/
