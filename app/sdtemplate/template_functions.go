package sdtemplate

import (
	"text/template"
	"strings"
)

var TemplateFuncMap = template.FuncMap{
	"escape_tag_key": EscapeTagKey,
}

func EscapeTagKey(key string) string {
	return strings.Replace(key, ".", "_", -1)
}
