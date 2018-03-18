package docker

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"github.com/davidhiendl/telegraf-docker-sd/app/backend"
	"github.com/davidhiendl/telegraf-docker-sd/app/utils"
)

const (
	NAME = "docker"
)

type DockerBackend struct {
	commonConfig     *config.ConfigSpec
	templates        map[string]*sdtemplate.Template
	telegrafReloader *utils.TelegrafReloader

	config            *DockerConfigSpec
	dockerCli         *client.Client
	dockerCtx         context.Context
	tags              []string
	trackedContainers map[string]*TrackedContainer
}

func NewBackend() *DockerBackend {
	return &DockerBackend{
		trackedContainers: make(map[string]*TrackedContainer),
	}
}

func (backend *DockerBackend) Name() string {
	return NAME
}

func (backend *DockerBackend) Status() int {
	return 1
}

func (backend *DockerBackend) Init(spec *backend.BackendConfigSpec) {
	backend.commonConfig = spec.Config
	backend.templates = spec.Templates
	backend.telegrafReloader = spec.Reloader
	backend.config = LoadConfig()
	backend.prepareDockerClient()
}

func (backend *DockerBackend) Run() {
	backend.processContainers()
}

func (backend *DockerBackend) Clean() {
	for id, container := range backend.trackedContainers {
		container.RemoveConfigFile()
		delete(backend.trackedContainers, id)
	}
	backend.telegrafReloader.ShouldReload = true
}
