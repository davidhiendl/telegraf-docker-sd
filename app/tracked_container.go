package app

import (
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"os"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	"path/filepath"
)

// TrackedContainer is used to maintain state about already processed containers and to be able to remove their configurations easily
type TrackedContainer struct {
	app         *App
	containerID string
	params      *sdtemplate.Params
	configFile  string
}

// Create new config and populate it from environment
func NewTrackedContainer(app *App, containerID string, params *sdtemplate.Params) (tracked *TrackedContainer) {
	c := TrackedContainer{
		app:         app,
		containerID: containerID,
		params:      params,
	}
	return &c
}

func (tc *TrackedContainer) GetConfigFile() string {
	if tc.configFile == "" {
		file, _ := filepath.Abs(tc.app.config.ConfigDir + "/" + tc.app.config.AutoConfPrefix + tc.GetShortID() + tc.app.config.AutoConfExtension)
		tc.configFile = file
	}
	return tc.configFile
}

func (tc *TrackedContainer) GetShortID() string {
	return tc.containerID[0:12]
}

func (tc *TrackedContainer) RemoveConfigFile() {
	stat, err := os.Stat(tc.GetConfigFile())
	if os.IsNotExist(err) {
		return
	}

	// do not touch anything that is not a file
	if stat.IsDir() {
		logger.Errorf("Config file is directory: %v \n", tc.GetConfigFile())
		return
	}

	err = os.Remove(tc.GetConfigFile())
	if err != nil {
		logger.Errorf("Failed to remove config file: %v with err: %v", tc.GetConfigFile(), err)
	}
}

func (tc *TrackedContainer) WriteConfigFile(contents string) {
	// open file using READ & WRITE permission
	file, err := os.OpenFile(tc.GetConfigFile(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Errorf("%v: failed to open file: %v err: %v", tc.GetShortID(), tc.GetConfigFile, err)
		panic(err)
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(contents)
	if err != nil {
		logger.Errorf("%v: failed to write file: %v err: %v", tc.GetShortID(), tc.GetConfigFile, err)
		panic(err)
	}

	// save changes
	err = file.Sync()
	if err != nil {
		logger.Errorf("%v: failed to sync file: %v err: %v", tc.GetShortID(), tc.GetConfigFile, err)
		panic(err)
	}

	logger.Infof("%v: wrote configuration: %v", tc.GetShortID(), tc.GetConfigFile())
}
