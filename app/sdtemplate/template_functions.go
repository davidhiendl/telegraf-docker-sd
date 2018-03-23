package sdtemplate

import (
	"text/template"
	"strings"
	"fmt"
	"strconv"
	"github.com/sirupsen/logrus"
)

var TemplateFuncMap = template.FuncMap{
	"escape_tag_key":   EscapeTagKey,
	"format_indent":    FormatIndent,
	"as_key_value_map": AsKeyValueMap,
	"to_toml_array":    ToTomlArray,
}

// escape a tag key for influxdb
func EscapeTagKey(key string) string {
	// TODO needs work/regex to strip everything
	return strings.Replace(key, ".", "_", -1)
}

// indent a group of values by a given amount of spaces
func FormatIndent(value string, spaces string) string {

	spacesInt, err := strconv.Atoi(spaces)
	if err != nil {
		logrus.Fatalf(`failed to parse indent: value="%v" count="%v"`, value, spaces)
	}

	rows := strings.Split(value, "\n")

	prefix := strings.Repeat(" ", spacesInt)

	for i := 0; i < len(rows); i++ {
		rows[i] = prefix + rows[i]
	}

	return strings.Join(rows, "\n")
}

func AsKeyValueMap(m map[string]string, indent int) string {
	rows := make([]string, len(m))

	prefix := strings.Repeat(" ", indent)

	i := 0
	for key, value := range m {
		rows[i] = prefix + fmt.Sprintf(`%v = "%v"`, key, value)
		i++
	}

	return strings.Join(rows, "\n")
}

func ToTomlArray(values []string) string {
	if len(values) <= 0 {
		return "[]"
	}

	// TODO escape values
	return `["` + strings.Join(values, `","`) + `"]`
}
