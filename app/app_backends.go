package app

import (
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"github.com/davidhiendl/telegraf-docker-sd/app/backend/docker"
	"github.com/davidhiendl/telegraf-docker-sd/app/backend"
	"github.com/davidhiendl/telegraf-docker-sd/app/backend/kubernetes"
	"github.com/sirupsen/logrus"
)

// Load templates from disk. If called multiple times templates are re-loaded
func (app *App) loadBackends() {

	for _, name := range app.config.Backends {
		// skip already loaded backends
		if _, ok := app.backends[name]; ok {
			continue
		}

		// create backend instance
		logrus.Debugf("[%v] creating backend", name)
		b := LoadBackend(name)

		// filter templates by backend
		templates := make(map[string]*sdtemplate.Template)
		for tplName, tpl := range app.templates {
			if tpl.Spec.Backend == b.Name() {
				logrus.Debugf("[%v] adding template: %v", name, tplName)
				templates[tplName] = tpl
			}
		}

		// prepare backend config
		spec := &backend.BackendConfigSpec{
			Config:    app.config,
			Reloader:  app.telegrafReloader,
			Templates: templates,
		}

		// configure backend
		b.Init(spec)

		app.backends[name] = b
	}
}

func LoadBackend(name string) backend.Backend {

	switch name {
	case docker.NAME:
		return docker.NewBackend()
	case kubernetes.NAME:
		return kubernetes.NewBackend()
	default:
		logrus.Fatalf(`unknown backend: "%v"`, name)
	}

	return nil
}
