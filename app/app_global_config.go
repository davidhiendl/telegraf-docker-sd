package app

import (
	"strings"
	"path/filepath"
	"bytes"
	"os"
	"github.com/davidhiendl/telegraf-docker-sd/app/globalconfig"
	"github.com/sirupsen/logrus"
)

func (app *App) processGlobalConfig() {

	globalConfig := &globalconfig.GlobalConfigSpec{
		EnvMap:   app.config.EnvMap,
		Tags:     app.config.GlobalTags,
		Backends: app.config.Backends,
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
			logrus.Fatalf("[global][%v] error during template execution: %+v", simpleName, err)
		}

		// write out config file
		writeConfigFile(configFile, configBuffer.String())
		logrus.Debugf("[global][%v] wrote main config file: %v", simpleName, configFile)
	}

}
func writeConfigFile(path string, contents string) {
	// open file using READ & WRITE permission
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logrus.Errorf("[global] failed to open file: %v err: %v", path, err)
		panic(err)
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(contents)
	if err != nil {
		logrus.Errorf("[global] failed to write file: %v err: %v", path, err)
		panic(err)
	}

	// save changes
	err = file.Sync()
	if err != nil {
		logrus.Errorf("[global] failed to sync file: %v err: %v", path, err)
		panic(err)
	}
}
