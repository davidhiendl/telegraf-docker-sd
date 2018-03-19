package docker

import "regexp"

func (tc *TrackedContainer) parseSwarmLabelsAsTags() {
	for _, swarmLabel := range SWARM_LABELS {
		for label, value := range tc.Container.Labels {
			if label == swarmLabel {
				tc.Tags[label] = value
				break
			}
		}
	}
}

func (tc *TrackedContainer) parseLabelsAsTags() {

	// explicit labels
	rex, err := regexp.Compile("^" + regexp.QuoteMeta("telegraf.sd.tags.") + "(.+)$")
	if err != nil {
		panic(err)
	}

	for key, value := range tc.Container.Labels {
		matches := rex.FindAllStringSubmatch(key, -1)
		if matches != nil {
			shortName := matches[0][1]
			tc.Tags[shortName] = value
			continue
		}
	}

	// whitelist mode
	if len(tc.backend.config.TagLabelsWhitelist) > 0 {
		for _, label := range tc.backend.config.TagLabelsWhitelist {
			value, ok := tc.Container.Labels[label]
			if ok {
				tc.Tags["label_"+label] = value
			}
		}
	}

	// blacklist mode
	if len(tc.backend.config.TagLabelsBlacklist) > 0 {
		for label, value := range tc.Container.Labels {

			// check if label is in blacklist
			found := false
			for _, blacklisted := range tc.backend.config.TagLabelsBlacklist {
				if label == blacklisted {
					found = true
					break
				}
			}

			// if label was not found in blacklist
			if !found {
				tc.Tags["label_"+label] = value
			}
		}
	}
}
