package sdtemplate

import (
	"text/template"
	"io"
	"path/filepath"
	"regexp"
)

type Template struct {
	template *template.Template;
	Spec     *TemplateSpec;
	FileName string;
}

// Create new config and populate it from environment
func NewTemplate(filePath string) (*Template, error) {

	// parse spec
	spec, err := ParseSpecFile(filePath)
	if err != nil {
		return nil, err
	}

	// get base filename as reference for the template
	fileName := filepath.Base(filePath)

	// create template
	goTpl := template.New(fileName)
	goTpl.Funcs(TemplateFuncMap)
	goTpl.Parse(spec.Template)

	// assemble object
	tpl := Template{
		Spec:     spec,
		template: goTpl,
		FileName: fileName,
	}

	return &tpl, err
}

func (tpl *Template) Execute(wr io.Writer, data interface{}) error {
	return tpl.template.Execute(wr, data)
}

// remove comments from templates
// TODO make use of this again
func (tpl *Template) cleanTemplateOutput(contents string) (string) {

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
