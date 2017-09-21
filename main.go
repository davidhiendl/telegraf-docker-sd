package main

import (
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"fmt"
	"github.com/davidhiendl/telegraf-docker-sd/app"
	"os"
)

func main() {
	// retrieve config
	config := app.NewConfig()
	err := config.LoadFromEnv()
	if err != nil {
		panic(err);
	}

	fmt.Printf("Config: %+v\n", config)

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	instance := app.NewApp(config, cli, ctx)

	switch {

	case len(os.Args) <= 0:
		instance.ClearConfigFiles()
		instance.Watch()

	case os.Args[0] == "run":
		instance.ClearConfigFiles()
		instance.Run()

	case os.Args[0] == "watch":
		instance.ClearConfigFiles()
		instance.Watch()

	case os.Args[0] == "clear":
		instance.ClearConfigFiles()
		instance.TriggerTelegrafReload()

	default:
		instance.ClearConfigFiles()
		instance.Watch()
	}
}
