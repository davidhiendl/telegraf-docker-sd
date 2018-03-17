package app

import (
	"strings"
	"path/filepath"
	"bytes"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	"os"
	"github.com/davidhiendl/telegraf-docker-sd/app/globalconfig"
)

func (app *App) processGlobalConfig() {

	globalConfig := &globalconfig.GlobalConfigSpec{
		EnvMap:     app.config.EnvMap,
		Tags: app.config.GlobalTags,
		Backends:   app.config.Backends,
	}

	for _, template := range app.templates {
		// filter by backend type
		if template.Spec.Backend != "global" {
			continue
		}

		simpleName := strings.TrimSuffix(strings.TrimSuffix(template.FileName, ".yaml"), ".yml")
		configFile, _ := filepath.Abs(
			app.config.ConfigDir +
				"/" + app.config.AutoConfPrefix +
				"_global_" +
				simpleName +
				app.config.AutoConfExtension)

		configBuffer := new(bytes.Buffer)
		err := template.Execute(configBuffer, globalConfig)
		if err != nil {
			logger.Fatalf("[global][%v] error during template execution: %+v", simpleName, err)
		}

		// write out config file
		writeConfigFile(configFile, configBuffer.String())
		logger.Debugf("[global][%v] wrote main config file: %v", configFile)
	}

}
func writeConfigFile(path string, contents string) {
	// open file using READ & WRITE permission
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Errorf("[global] failed to open file: %v err: %v", path, err)
		panic(err)
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(contents)
	if err != nil {
		logger.Errorf("[global] failed to write file: %v err: %v", path, err)
		panic(err)
	}

	// save changes
	err = file.Sync()
	if err != nil {
		logger.Errorf("[global] failed to sync file: %v err: %v", path, err)
		panic(err)
	}

	logger.Infof("[global] wrote configuration: %v", path, err)
}
