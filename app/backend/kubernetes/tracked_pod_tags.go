package kubernetes

import (
	"regexp"
)

func (tp *TrackedPod) parseSystemTags() {

	if tp.backend.config.TagNamespace {
		tp.Tags["k8s_namespace"] = tp.Pod.Namespace
	}

	if tp.backend.config.TagPodName {
		tp.Tags["k8s_pod"] = tp.Pod.Name
	}
}

func (tp *TrackedPod) parseLabelsAsTags() {

	// whitelist mode
	if len(tp.backend.config.TagLabelsWhitelist) > 0 {
		for _, label := range tp.backend.config.TagLabelsWhitelist {
			value, ok := tp.Pod.Labels[label]
			if ok {
				tp.Tags["k8s_label_"+label] = value
			}
		}
	}

	// blacklist mode
	if len(tp.backend.config.TagLabelsBlacklist) > 0 {
		for label, value := range tp.Pod.Labels {

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
				tp.Tags["k8s_label_"+label] = value
			}
		}
	}
}

// extract configuration values from annoations
func (tp *TrackedPod) parseAnnotationsAsTags() {
	rex, err := regexp.Compile("^" + regexp.QuoteMeta("telegraf.sd.tags/") + "(.+)$")
	if err != nil {
		panic(err)
	}

	for key, value := range tp.Pod.Annotations {
		matches := rex.FindAllStringSubmatch(key, -1)
		if matches != nil {
			shortName := matches[0][1]
			tp.Tags[shortName] = value
			continue
		}
	}
}
