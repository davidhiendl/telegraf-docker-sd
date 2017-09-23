package app

import (
	"github.com/davidhiendl/telegraf-docker-sd/app/sdtemplate"
	"io/ioutil"
	"log"
	"regexp"
	"github.com/davidhiendl/telegraf-docker-sd/app/logger"
	"github.com/davidhiendl/telegraf-docker-sd/app/tgtemplate"
	"os"
	"bytes"
)

const TELEGRAF_MAIN_TEMPLATE_SRC_FILE = "_telegraf.goconf"
const TELEGRAF_MAIN_TEMPLATE_DST_FILE = "telegraf.conf"

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

		// load telegraf main template file (for telegraf.conf)
		if f.Name() == TELEGRAF_MAIN_TEMPLATE_SRC_FILE {
			tpl, err := tgtemplate.NewTemplate(name, filePath)
			if err != nil {
				panic(err)
			}
			app.telegrafTemplate = tpl
		} else { // or load regular template file
			tpl, err := sdtemplate.NewTemplate(name, filePath)
			if err != nil {
				panic(err)
			}
			app.templates[name] = tpl
		}
	}

}

// process main template file if it exists
func (app *App) processMainTemplateFile() (bool) {
	if app.telegrafTemplate == nil {
		return false
	}

	configBuffer := new(bytes.Buffer)
	err := app.telegrafTemplate.Execute(tgtemplate.NewParams(), configBuffer)
	if err != nil {
		logger.Errorf("Failed to process main config file")
		panic(err)
	}

	app.writeMainConfigFile(app.cleanTemplateOutput(configBuffer.String()))
	logger.Infof("Wrote main configuration: %v", app.config.ConfigDir+"/"+TELEGRAF_MAIN_TEMPLATE_DST_FILE)

	return true
}

func (app *App) writeMainConfigFile(contents string) {
	// open file using READ & WRITE permission
	target := app.config.ConfigDir + "/" + TELEGRAF_MAIN_TEMPLATE_DST_FILE
	file, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		logger.Errorf("Failed to write main config file: %v", target)
		panic(err)
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(contents)
	if err != nil {
		logger.Errorf("Failed to write main config file: %v", target)
		panic(err)
	}

	// save changes
	err = file.Sync()
	if err != nil {
		logger.Errorf("Failed to write main config file: %v", target)
		panic(err)
	}
}

func (app *App) cleanTemplateOutput(contents string) (string) {
	if !app.config.CleanOutput {
		return contents
	}


	regexCleanComments, err := regexp.Compile("(?m:^\\s*#.*$\n)")
	if err != nil {
		panic(err)
	}

	/*
	regexCleanEmptyLines, err := regexp.Compile("(?m:^\\s*$\\n)")
	if err != nil {
		panic(err)
	}
	*/

	return regexCleanComments.ReplaceAllString(contents, "")
	// return regexCleanEmptyLines.ReplaceAllString(regexCleanComments.ReplaceAllString(contents, ""), "\n")
}
