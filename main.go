package main

import (
	"github.com/davidhiendl/telegraf-docker-sd/app"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
	"github.com/fatih/structs"
)

func main() {

	cfg := config.Load()
	logger.SetLevel(logger.LOG_DEBUG)

	// print config
	m := structs.Map(cfg)
	logger.Infof("[global] configuration loaded:")
	for key, value := range m {
		logger.Infof("%v = %v", key, value)
	}

	instance := app.NewApp(cfg)
	instance.ClearConfigFiles()
	instance.Run()
}
