package docker

import (
	"github.com/docker/docker/api/types"
	"context"
	"github.com/docker/docker/client"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

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
	cli, err := client.NewEnvClient()
	if err != nil {
		logger.Fatalf(LOG_PREFIX+" failed to connect to docker: %+v\n", err)
	}

	// allow selecting API version dynamically
	cli.NegotiateAPIVersion(ctx)

	backend.dockerCtx = ctx
	backend.dockerCli = cli
}
