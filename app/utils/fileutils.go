package utils

import (
	"os"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func RemoveConfigFile(path string) {
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return
	}

	// do not touch anything that is not a file
	if stat.IsDir() {
		logger.Errorf("[template] config file is directory: %v", path)
		return
	}

	err = os.Remove(path)
	if err != nil {
		logger.Errorf("[template] failed to remove config file: %v with err: %v", path, err)
	}
}

func WriteConfigFile(path string, contents string) {
	// open file using READ & WRITE permission
	fd, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		logger.Errorf("[template] failed to open file: %v err: %v", path, err)
		panic(err)
	}
	defer fd.Close()

	// write some text line-by-line to file
	_, err = fd.WriteString(contents)
	if err != nil {
		logger.Errorf("[template] failed to write file: %v err: %v", path, err)
		panic(err)
	}

	// save changes
	err = fd.Sync()
	if err != nil {
		logger.Errorf("[template] failed to sync file: %v err: %v", path, err)
		panic(err)
	}
}
