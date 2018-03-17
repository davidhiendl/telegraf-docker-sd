package docker

import (
	"docker.io/go-docker/api/types"
	"strings"
	"context"
	"docker.io/go-docker"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (backend *DockerBackend) processConfig() {
	backend.processConfigLabelsAsTags()
}

func (backend *DockerBackend) processConfigLabelsAsTags() {
	labelsRaw := strings.Split(backend.config.TagsFromLabels, ",")

	labelsClean := make([]string, len(labelsRaw))
	i := 0
	for _, label := range labelsRaw {
		labelsClean[i] = label
		i++
	}

	backend.tags = labelsClean
}

func (backend *DockerBackend) getImageForID(id string) *types.ImageSummary {
	images, err := backend.dockerCli.ImageList(backend.dockerCtx, types.ImageListOptions{})
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

func (backend *DockerBackend) prepareDockerClient() {
	ctx := context.Background()
	cli, err := docker.NewEnvClient()
	if err != nil {
		logger.Fatalf("[docker] failed to connect to docker: %+v\n", err)
	}

	backend.dockerCtx = ctx
	backend.dockerCli = cli
}
