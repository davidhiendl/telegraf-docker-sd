package docker

import (
	"github.com/docker/docker/api/types"
	"context"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
)

func (backend *DockerBackend) getImageForID(id string) (*types.ImageSummary, error) {
	images, err := backend.dockerCli.ImageList(backend.dockerCtx, types.ImageListOptions{})
	if err != nil {
		return nil, err
	}

	for _, image := range images {
		if image.ID == id {
			return &image, nil
		}
	}

	return nil, nil
}

func (backend *DockerBackend) prepareDockerClient() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		logrus.Errorf(LOG_PREFIX+" failed to connect to docker: %+v\n", err)
		return
	}

	// allow selecting API version dynamically
	cli.NegotiateAPIVersion(ctx)

	backend.dockerCtx = ctx
	backend.dockerCli = cli
}

func (backend *DockerBackend) resetDockerClient() {
	if backend.dockerCli != nil {
		backend.dockerCli.Close()
		backend.dockerCli = nil
	}

	if backend.dockerCtx != nil {
		backend.dockerCtx = nil
	}
}

func toShortID(id string) string {
	return id[0:12]
}
