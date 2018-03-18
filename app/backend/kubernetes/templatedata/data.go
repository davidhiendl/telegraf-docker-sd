package templatedata

type TemplateData struct {
	Name      string
	EnvMap    map[string]string
	Tags      map[string]string
	Config    map[string]string
}

func NewTemplateData() *TemplateData {
	td := &TemplateData{
		Tags:      make(map[string]string),
	}

	// td.parseAnnotations()

	return td
}


func (td *TemplateData) TargetIP() string {
	return ""
}
