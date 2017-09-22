# Templates
Automatically generating the main telegraf configuratio is also supported via a special template `_telegraf.goconf` in the template folder.
This template is only generated once at startup.

## Available functions and variables

| Name               | Type     | Params                     | Return              | Description                                                                |
| ---                | ---      | ---                        | ---                 | ---                                                                        |
| .EnvOrDefault      | Function | label string, value string | string              | Get a environment variable or the default                                  |
| .EnvGet            | Function | label string               | string              | Get a environment variable                                                 |
| .EnvHas            | Function | label string               | bool                | Check if a environment variable exists                                     |
| .EnvEquals         | Function | label string, value string | bool                | Check if a environment variables equals a given value                      |
| .GlobalTagFromEnv  | Function | label string               | string              | Generate a tag string from a given environment variable or an empty string |
| .GlobalTagsFromEnv | Function | none                       | map\[string\]string | Get a key => value map of tags                                             |


## Examples

**Check environment variables exist variables:**
```
{{ if .EnvHas "SOME_VAR" }}
    {{ .Env "SOME_VAR" }}
{{end}}
```

**Use environment or default**
```
{{ .EnvOrDefault "SOME_VAR" "default-value" }} # outputs the value of SOME_VAR or the default
```

**Apply global tags tags:**
```
# Global tags can be specified here in key="value" format.
[global_tags]
  # import all environment variables with format "GLOBAL_TAGS_$key=$value" as tags
  {{ range $key, $value := .GlobalTagsFromEnv }}
  {{ $key }} = "{{ $value }}"
  {{ end }}
```
