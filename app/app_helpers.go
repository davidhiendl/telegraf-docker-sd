package app

import (
	"strings"
	"github.com/docker/docker/api/types"
)

func (app *App) processConfig() {
	app.processConfigLabelsAsTags()
}

func (app *App) processConfigLabelsAsTags() {
	labelsRaw := strings.Split(app.config.TagsFromLabels, ",")

	labelsClean := []string{}
	for _, label := range labelsRaw {
		labelsClean = append(labelsClean, label)
	}

	app.tagsFromLabels = labelsClean
}

func (app *App) cleanupTrackedContainer(tracked *TrackedContainer) {
	tracked.RemoveConfigFile()
	delete(app.trackedContainers, tracked.containerID)
	app.shouldReload = true
}

func (app *App) getImageForID(id string) *types.ImageSummary {
	images, err := app.docker.ImageList(app.ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, image := range images {
		if image.ID == id {
			return &image
		}
	}

	return nil
}
