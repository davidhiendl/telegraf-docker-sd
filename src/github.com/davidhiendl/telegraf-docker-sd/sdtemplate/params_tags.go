package sdtemplate

import (
	"regexp"
)

// extract configuration values from labels
func (params *Params) parseExplicitLabelsAsTags() {
	rex, err := regexp.Compile("^telegraf\\.sd\\.tags\\.([a-zA-Z0-9_\\.\\-]*)$")
	if err != nil {
		panic(err)
	}

	for label, value := range params.Container.Labels {

		// add explicit labels as tags
		matches := rex.FindAllStringSubmatch(label, -1)
		if matches != nil {
			shortName := matches[0][1]
			params.Tags[shortName] = value
			continue
		}
	}
}

// check and import swarm labels
func (params *Params) ParseLabelsAsTags(labelsAsTags []string) {
	for label, value := range params.Container.Labels {
		for _, matchLabel := range labelsAsTags {
			if label == matchLabel {
				params.Tags[label] = value
				break
			}
		}
	}
}
