package app

import (
	"time"
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
	"github.com/davidhiendl/telegraf-docker-sd/app/utils"
	"github.com/davidhiendl/telegraf-docker-sd/app/backend"
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"github.com/davidhiendl/telegraf-docker-sd/app/constants"
	"regexp"
	"os"
	"io/ioutil"
	"github.com/sirupsen/logrus"
)

type App struct {
	config                *config.ConfigSpec
	telegrafReloader      *utils.TelegrafReloader
	processedMainTemplate bool
	backends              map[string]backend.Backend
	templates             map[string]*sdtemplate.Template
}

// Create new config and populate it from environment
func NewApp(cfg *config.ConfigSpec) (*App) {
	app := App{
		config:                cfg,
		processedMainTemplate: false,
		backends:              make(map[string]backend.Backend),
	}

	// register telegraf reload handler
	app.telegrafReloader = utils.NewTelegrafReloader()
	logrus.Infof("created reloader: %+v", app.telegrafReloader)

	app.loadTemplates()
	app.loadBackends()

	return &app
}

func (app *App) Run() {
	interval := app.config.QueryInterval
	if interval <= 0 {
		interval = constants.DEFAULT_QUERY_INTERVAL
	}
	dInterval := time.Duration(interval) * time.Second

	app.processGlobalConfig()

	for {
		for _, b := range app.backends {
			logrus.Debugf("run backend: %v", b.Name())
			b.Run()
		}
		app.telegrafReloader.ReloadIfRequested()
		time.Sleep(dInterval)
	}
}

// remove all configuration files that were created by regex: starting with prefix and ending with extension
func (app *App) ClearConfigFiles() {
	files, err := ioutil.ReadDir(app.config.ConfigDir)
	if err != nil {
		logrus.Fatalf("failed to clear config: %v", err)
	}

	// summarized: ^prefix[a-zA-Z0-9._-]*extension$
	rex, _ := regexp.Compile(
		"^" +
			regexp.QuoteMeta(app.config.AutoConfPrefix) +
			"[a-zA-Z0-9_\\.\\-]*" +
			regexp.QuoteMeta(app.config.AutoConfExtension) +
			"$")

	for _, f := range files {
		if rex.MatchString(f.Name()) {
			logrus.Debugf("cleaning up file: %v", f.Name())
			path := app.config.ConfigDir + "/" + f.Name()

			stat, err := os.Stat(path)
			if os.IsNotExist(err) {
				continue
			}

			// do not touch anything that is not a file
			if !stat.Mode().IsRegular() {
				logrus.Debugf("Config file is not a regular file: %v", path)
				continue
			}

			err = os.Remove(path)
			if err != nil {
				logrus.Debugf("failed to remove file: %v, err: %v", path, err)
			}
		}
	}
}
