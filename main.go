package main

import (
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"os"
	"github.com/davidhiendl/telegraf-docker-sd/app"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

func main() {

	// retrieve config
	config := app.NewConfig()
	err := config.LoadFromEnv()
	if err != nil {
		logger.Fatalf("failed to parse configuration from environment: %v \n", err)
	}

	logger.SetLevel(config.LogLevel)

	// print config
	m := config.AsMap()
	for key, value := range m {
		logger.Infof("Config.%v = %v", key, value)
	}

	// create docker connection
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		logger.Fatalf("Config: %+v\n", config)
	}

	instance := app.NewApp(config, cli, ctx)

	switch {

	case len(os.Args) <= 0:
		logger.Infof("starting to watch containers")
		instance.ClearConfigFiles()
		instance.Watch()

	case os.Args[0] == "run":
		logger.Infof("generating configuration for containers")
		instance.ClearConfigFiles()
		instance.Run()

	case os.Args[0] == "watch":
		logger.Infof("starting to watch containers")
		instance.ClearConfigFiles()
		instance.Watch()

	case os.Args[0] == "clear":
		logger.Infof("clearing existing auto-generated configuration files")
		instance.ClearConfigFiles()
		instance.Reload()

	default:
		logger.Infof("starting to watch containers")
		instance.ClearConfigFiles()
		instance.Watch()
	}
}
