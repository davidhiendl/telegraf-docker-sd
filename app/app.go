package app

import (
	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"golang.org/x/net/context"
	"log"
	"io/ioutil"
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"bytes"
	"time"
	"os"
	"syscall"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
)

type App struct {
	config            *Config
	tagsFromLabels    []string
	docker            *client.Client
	ctx               context.Context
	templates         map[string]*sdtemplate.Template
	trackedContainers map[string]*TrackedContainer
	shouldReload      bool
	signalDispatcher  []*SignalDispatcher
}

// Create new config and populate it from environment
func NewApp(config *Config, cli *client.Client, ctx context.Context) (*App) {
	app := App{
		config:            config,
		docker:            cli,
		ctx:               ctx,
		trackedContainers: make(map[string]*TrackedContainer),
		shouldReload:      false,
	}

	app.processConfig()
	app.loadTemplates()

	// register telegraf reload handler
	app.signalDispatcher = append(app.signalDispatcher, NewSignalHandler("telegraf", syscall.SIGHUP))

	return &app
}

// periodically execute Run
func (app *App) Watch() {
	raw := app.config.QueryInterval
	if raw <= 0 {
		raw = CONFIG_DEFAULT_QUERY_INTERVAL
	}

	interval := time.Duration(raw) * time.Second

	for {
		app.Run()
		time.Sleep(interval)
	}
}

// run templates against containers and generate config
func (app *App) Run() {
	app.ProcessContainers()
	if app.shouldReload {
		app.Reload()
	}
}

// reloads all registered agents
func (app *App) Reload() {
	// fmt.Println("reloading configuration")

	for _, sh := range app.signalDispatcher {
		sh.Execute()
	}

	app.shouldReload = false
}

// remove all configuration files that were created by regex: starting with prefix and ending with extension
func (app *App) ClearConfigFiles() {
	files, err := ioutil.ReadDir(app.config.ConfigDir)
	if err != nil {
		log.Fatal(err)
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
			logger.Debugf("cleaning up file: %v", f.Name())
			path := app.config.ConfigDir + "/" + f.Name()

			stat, err := os.Stat(path)
			if os.IsNotExist(err) {
				continue
			}

			// do not touch anything that is not a file
			if stat.IsDir() {
				logger.Debugf("Config file is directory: %v", path)
				continue
			}

			err = os.Remove(path)
			if err != nil {
				panic(err)
			}
		}
	}
}

// Load templates from disk. If called multiple times templates are re-loaded
func (app *App) loadTemplates() {
	app.templates = make(map[string]*sdtemplate.Template)

	files, err := ioutil.ReadDir(app.config.TemplateDir)
	if err != nil {
		log.Fatal(err)
	}

	// check all files and extract simple name without extension
	rex, _ := regexp.Compile("(^[a-zA-Z0-9_\\.\\-]*)\\.goconf$")
	for _, f := range files {
		matches := rex.FindAllStringSubmatch(f.Name(), -1)
		if matches == nil {
			continue
		}

		name := matches[0][1]
		filePath := app.config.TemplateDir + "/" + f.Name()
		logger.Infof("loading config template: %v from %v", name, filePath)

		tpl, err := sdtemplate.NewTemplate(name, filePath)
		if err != nil {
			panic(err)
		}
		app.templates[name] = tpl
	}

}

func (app *App) ProcessContainers() {
	containers, err := app.docker.ContainerList(app.ctx, types.ContainerListOptions{});
	if err != nil {
		logger.Warnf("failed to process container: %v", err)
		return
	}

	// check existing containers and configure them
	for _, cont := range containers {
		app.ProcessContainer(cont)
	}

	// iterate over all currently tracked containers and clean up their config files
	for id, tracked := range app.trackedContainers {
		found := false

		// check if it still exists
		for _, cont := range containers {
			if cont.ID == id {
				found = true
			}
		}

		// if it does not exist anymore then remove the associated config
		if !found {
			logger.Debugf("cleaning up no longer tracked container: %v", id)
			app.cleanupTrackedContainer(tracked)
		}
	}
}

func (app *App) ProcessContainer(cont types.Container) {

	// check if running
	running := cont.State == "running"

	// check if container already tracked
	if tracked, ok := app.trackedContainers[cont.ID]; ok {
		if !running {
			app.cleanupTrackedContainer(tracked)
		}
		return
	}

	// do not configure if not running
	if !running {
		return
	}

	// check if bridge network exists
	_, ok := cont.NetworkSettings.Networks["bridge"]
	if !ok {
		logger.Debugf("%v: missing network bridge on container, skipping", cont.Names[0])
		return
	}

	// assemble template params
	image := app.getImageForID(cont.ImageID)
	params := sdtemplate.NewParams(cont, image)

	// add swarm labels if desired
	if app.config.TagsFromSwarmLabels {
		params.ParseLabelsAsTags(SWARM_LABELS)
	}

	logger.Debugf("%v: detected tags: %+v", cont.Names[0], params.Tags)
	logger.Debugf("%v: detected config: %+v", cont.Names[0], params.Config)

	// register tracked container
	tracked := NewTrackedContainer(app, cont.ID, params)
	app.trackedContainers[cont.ID] = tracked

	// process template(s) for container
	confStr := app.processTemplatesAgainstContainer(params)
	tracked.WriteConfigFile(confStr)

	// mark as changed
	app.shouldReload = true
}

func (app *App) processTemplatesAgainstContainer(params *sdtemplate.Params) string {
	buf := new(bytes.Buffer)
	for _, template := range app.templates {
		err := template.Execute(params, buf)
		if err != nil {
			panic(err)
		}
	}

	return buf.String()
}
