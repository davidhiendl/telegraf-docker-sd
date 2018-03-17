package sdtemplate

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
)

func ParseSpec(raw []byte) (*TemplateSpec, error) {
	spec := &TemplateSpec{}
	err := yaml.Unmarshal(raw, spec)
	if err != nil {
		return nil, err
	}

	return spec, nil
}

func ParseSpecFile(file string) (*TemplateSpec, error) {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	return ParseSpec(raw)
}

type TemplateSpec struct {
	Name         string        `yaml:"name"`
	Backend      string        `yaml:"backend"`
	Template     string        `yaml:"template"`
	TemplateType string        `yaml:"templateType"`
	Labels       []string      `yaml:"labels"`
	Matchers     []MatcherSpec `yaml:"matchers"`
}

type MatcherSpec struct {
	// TODO implement yaml based filtering instead of go template
}
