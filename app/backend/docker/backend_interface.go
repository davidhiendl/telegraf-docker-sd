package docker

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"github.com/davidhiendl/telegraf-docker-sd/app/backend"
	"github.com/davidhiendl/telegraf-docker-sd/app/utils"
	"github.com/sirupsen/logrus"
	"github.com/fatih/structs"
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

	// print config
	m := structs.Map(backend.config)
	for key, value := range m {
		logrus.WithFields(logrus.Fields{"key": key, "value": value}).Infof(LOG_PREFIX + " configuration loaded")
	}

	backend.prepareDockerClient()
}

func (backend *DockerBackend) Run() {

	// check if client exists and attempt re-connect if necessary
	if backend.dockerCli == nil {
		logrus.Errorf(LOG_PREFIX + " docker client not connected, attempting to reconnect")
		backend.prepareDockerClient()

		if backend.dockerCli == nil {
			logrus.Errorf(LOG_PREFIX + " failed to reconnect docker client")
			return
		}
	}

	backend.processContainers()
}

func (backend *DockerBackend) Clean() {
	for id, tc := range backend.trackedContainers {
		utils.RemoveConfigFile(tc.GetConfigFile())
		delete(backend.trackedContainers, id)
	}
	backend.telegrafReloader.RequestReload()
}
