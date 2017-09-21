package app

import (
	"github.com/davidhiendl/telegraf-docker-sd/sdtemplate"
	"os"
	"fmt"
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
		tc.configFile = tc.app.config.ConfigDir + "/" + tc.app.config.AutoConfPrefix + tc.GetShortID() + tc.app.config.AutoConfExtension
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
		fmt.Printf("ERROR: Config file is directory: %v \n", tc.GetConfigFile())
		return
	}

	err = os.Remove(tc.GetConfigFile())
	if err != nil {
		panic(err)
	}
}

func (tc *TrackedContainer) WriteConfigFile(contents string) {
	// open file using READ & WRITE permission
	file, err := os.OpenFile(tc.GetConfigFile(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(contents)
	if err != nil {
		panic(err)
	}

	// save changes
	err = file.Sync()
	if err != nil {
		panic(err)
	}
}
