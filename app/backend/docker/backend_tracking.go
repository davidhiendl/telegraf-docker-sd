package docker

import (
	"docker.io/go-docker/api/types"
	"bytes"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func (backend *DockerBackend) cleanupContainer(tracked *TrackedContainer) {
	logger.Debugf("[docker][%v] cleaning up no longer tracked container, file=%v", tracked.ShortID, tracked.GetConfigFile())
	tracked.RemoveConfigFile()
	delete(backend.trackedContainers, tracked.ID)
	backend.telegrafReloader.ShouldReload = true
}

func (backend *DockerBackend) processContainers() {

	containers, err := backend.dockerCli.ContainerList(backend.dockerCtx, types.ContainerListOptions{});
	if err != nil {
		logger.Warnf("[docker] failed to process containers: %v", err)
		return
	}

	// check existing containers and configure them
	for _, cont := range containers {
		backend.trackContainer(&cont)
	}

	// iterate over all currently tracked containers and clean up their config files
	for id, tracked := range backend.trackedContainers {
		found := false

		// check if it still exists
		for _, cont := range containers {
			if cont.ID == id {
				found = true
				break
			}
		}

		// if it does not exist anymore then remove the associated config
		if !found {
			logger.Debugf("[docker][%v] cleanup up container because it was no longer found", tracked.ShortID)
			backend.cleanupContainer(tracked)
		}
	}
}

func toShortID(id string) string {
	return id[0:12]
}

func (backend *DockerBackend) trackContainer(cont *types.Container) {

	// check if running
	running := cont.State == "running"

	// check if container already tracked
	if tracked, ok := backend.trackedContainers[cont.ID]; ok {
		// cleanup container that stopped running
		if !running {
			logger.Debugf("[docker][%v] cleanup up container because it is no longer running")
			backend.cleanupContainer(tracked)
		}
		return
	}

	// do not configure if not running
	if !running {
		return
	}

	// check if bridge network exists
	_, ok := cont.NetworkSettings.Networks["bridge"]
	if !ok {
		logger.Debugf("[docker][%v] missing network bridge on container, skipping", cont.Names[0])
		return
	}

	// register tracked container
	logger.Infof("[docker][%v] started tracking: %+v", toShortID(cont.ID), cont.Names)
	tracked := NewTrackedContainer(backend, cont)
	backend.trackedContainers[tracked.ID] = tracked

	// process template(s) for container
	configBuffer := new(bytes.Buffer)
	for _, template := range backend.templates {
		logger.Debugf("[docker][%v] running against template: %v", tracked.ShortID, template.FileName)
		err := template.Execute(configBuffer, tracked.Data)
		if err != nil {
			logger.Fatalf("[docker][%v] error during template execution: %+v", cont.Names[0], err)
		}
	}

	// write config
	tracked.WriteConfigFile(configBuffer.String())

	// mark as changed
	backend.telegrafReloader.ShouldReload = true
}


