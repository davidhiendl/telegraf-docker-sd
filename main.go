package main

import (
	"github.com/davidhiendl/telegraf-docker-sd/app"
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
	"github.com/fatih/structs"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {

	// configure logger
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	// load config
	cfg := config.Load()

	// print config
	m := structs.Map(cfg)
	for key, value := range m {
		if key == "EnvMap" {
			continue
		}
		logrus.WithFields(logrus.Fields{"key": key, "value": value}).Infof("configuration loaded")
	}

	instance := app.NewApp(cfg)
	instance.ClearConfigFiles()
	instance.Run()
}
