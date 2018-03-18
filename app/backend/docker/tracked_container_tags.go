package docker

import "regexp"

func (tp *TrackedContainer) parseSwarmLabelsAsTags() {
	for _, swarmLabel := range SWARM_LABELS {
		for label,value := range tp.Container.Labels {
			if label == swarmLabel {
				tp.Tags[label] = value
				break
			}
		}
	}
}

func (tp *TrackedContainer) parseLabelsAsTags() {

	// explicit labels
	rex, err := regexp.Compile("^" + regexp.QuoteMeta("telegraf.sd.tags.") + "(.+)$")
	if err != nil {
		panic(err)
	}

	for key, value := range tp.Container.Labels {
		matches := rex.FindAllStringSubmatch(key, -1)
		if matches != nil {
			shortName := matches[0][1]
			tp.Tags[shortName] = value
			continue
		}
	}

	// whitelist mode
	if len(tp.backend.config.TagLabelsWhitelist) > 0 {
		for _, label := range tp.backend.config.TagLabelsWhitelist {
			value, ok := tp.Container.Labels[label]
			if ok {
				tp.Tags["label_"+label] = value
			}
		}
	}

	// blacklist mode
	if len(tp.backend.config.TagLabelsBlacklist) > 0 {
		for label, value := range tp.Container.Labels {

			// check if label is in blacklist
			found := false
			for _, blacklisted := range tp.backend.config.TagLabelsBlacklist {
				if label == blacklisted {
					found = true
					break
				}
			}

			// if label was not found in blacklist
			if !found {
				tp.Tags["label_"+label] = value
			}
		}
	}
}

