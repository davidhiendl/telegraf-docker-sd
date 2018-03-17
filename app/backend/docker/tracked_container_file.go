package docker

import (
	"os"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	"path/filepath"
)

func (tc *TrackedContainer) GetConfigFile() string {
	if tc.configFile == "" {
		file, _ := filepath.Abs(
			tc.backend.commonConfig.ConfigDir +
				"/" + tc.backend.commonConfig.AutoConfPrefix +
				tc.backend.Name() + "_" +
				tc.ShortID +
				tc.backend.commonConfig.AutoConfExtension)
		tc.configFile = file
	}
	return tc.configFile
}

func (tc *TrackedContainer) RemoveConfigFile() {
	stat, err := os.Stat(tc.GetConfigFile())
	if os.IsNotExist(err) {
		return
	}

	// do not touch anything that is not a file
	if stat.IsDir() {
		logger.Errorf("[docker][%v] config file is directory: %v", tc.ShortID, tc.GetConfigFile())
		return
	}

	err = os.Remove(tc.GetConfigFile())
	if err != nil {
		logger.Errorf("[docker][%v] failed to remove config file: %v with err: %v", tc.ShortID, tc.GetConfigFile(), err)
	}
}

func (tc *TrackedContainer) WriteConfigFile(contents string) {
	// open file using READ & WRITE permission
	file, err := os.OpenFile(tc.GetConfigFile(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Errorf("[docker][%v] failed to open file: %v err: %v", tc.ShortID, tc.GetConfigFile(), err)
		panic(err)
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(contents)
	if err != nil {
		logger.Errorf("[docker][%v] failed to write file: %v err: %v", tc.ShortID, tc.GetConfigFile(), err)
		panic(err)
	}

	// save changes
	err = file.Sync()
	if err != nil {
		logger.Errorf("[docker][%v] failed to sync file: %v err: %v", tc.ShortID, tc.GetConfigFile(), err)
		panic(err)
	}

	logger.Infof("[docker][%v] wrote configuration: %v", tc.ShortID, tc.GetConfigFile())
}
