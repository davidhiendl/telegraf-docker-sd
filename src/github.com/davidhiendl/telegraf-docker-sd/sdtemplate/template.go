package sdtemplate

import (
	"text/template"
	"io"
	"path/filepath"
)

type Template struct {
	name     string;
	filePath string;
	fileName string;
	template *template.Template;
}

// Create new config and populate it from environment
func NewTemplate(name string, filePath string) (*Template, error) {
	tplw := Template{
		name:     name,
		filePath: filePath,
		fileName: filepath.Base(filePath),
	}

	gotpl, err := template.ParseFiles(tplw.filePath)
	tplw.template = gotpl

	if err != nil {
		panic(err)
	}

	return &tplw, err
}

func (tplw *Template) Name() string {
	return tplw.name
}

func (tplw *Template) Execute(params *Params, wr io.Writer) error {
	return tplw.template.ExecuteTemplate(wr, tplw.fileName, params)
}
