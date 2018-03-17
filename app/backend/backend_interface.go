package backend

import (
	"github.com/davidhiendl/telegraf-docker-sd/app/config"
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"github.com/davidhiendl/telegraf-docker-sd/app/utils"
)

type BackendConfigSpec struct {
	Config    *config.ConfigSpec
	Templates map[string]*sdtemplate.Template
	Reloader  *utils.TelegrafReloader
}

type Backend interface {
	Name() string
	Status() int
	Init(spec *BackendConfigSpec)
	Run()
	Clean()
}
