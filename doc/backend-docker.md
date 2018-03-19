# Backend: Docker

The discovery routine automatically scans the docker host for all of it's running containers and applies each template against the container.
It is the templates responsibility to determine if the container matches the template. This allows for maximum flexibility.


### Automatically collected tags and configurations
- Any container labels starting with `telegraf.sd.tags.*` will be added as to the tag list for this container.
- Any container labels starting with `telegraf.sd.config.*` will be shortened by this prefix and added to the config map.


## Environment Configuration Parameters

| Name                            | Type   | Default | Description                                                                                                                                                                |
| ---                             | ---    | ---     | ---                                                                                                                                                                        |
| TSD_DOCKER_TAGS_FROM_SWARM      | bool   | true    | Include swarm system labels as tags (com.docker.swarm.*, com.docker.stack.namespace)                                                                                                                                       |
| TSD_DOCKER_AUTO_CONF_PREFIX     | string | docker_ | Backend specific prefix for generated files                                                                                                                                |
| TSD_DOCKER_TAG_LABELS_WHITELIST | string |         | An explicit comma-separated list of labels to include as tags                                                                                                              |
| TSD_DOCKER_TAG_LABELS_BLACKLIST | string |         | A comma-separated list of labels to exclude from tags, if this value is non-empty all labels except the listed labels will be included. This may cause cardinality issues. |


## Available functions and variables


| Name              | Type     | Params                     | Return              | Description                                                  |
| ---               | ---      | ---                        | ---                 | ---                                                          |
| .BridgeIP         | Function | none                       | string              | The actual container bridge ip                               |
| .Tags             | Variable | none                       | map\[string\]string | A map of all computed tags and values                        |
| .MatchImage       | Function | pattern string             | bool                | Check container image against provided literal string        |
| .MatchImageRegex  | Function | key string                 | bool                | Check container image against provided regex                 |
| ---               | ---      | ---                        | ---                 | ---                                                          |
| .Config           | Variable | none                       | map\[string\]string | A map of all configurations values                           |
| .ConfigGet        | Function | key string                 | string              | Retrieve a config value from combined sources or use default |
| .ConfigOrDefault  | Function | key string, default string | string              | Retrieve a config value from combined sources or use default |
| .ConfigExists     | Function | key string                 | bool                | Check if a config key exists                                 |
| .ConfigEquals     | Function | key string, default string | bool                | Check if a config value equals an other value                |
| .ConfigMatches    | Function | key string, pattern string | bool                | Check if a config value matches a regex pattern              |
| ---               | ---      | ---                        | ---                 | ---                                                          |
| .Labels           | Variable | none                       | map\[string\]string | A map of all labels and their values                         |
| .LabelGet         | Function | key string                 | bool                | Get label value                                              |
| .LabelOrDefault   | Function | key string, default string | bool                | Get label value or default                                   |
| .LabelExists      | Function | key string                 | bool                | Check if label exists                                        |
| .LabelEquals      | Function | key string, compare string | bool                | Check if label equals value                                  |
| .LabelMatches     | Function | key string, pattern string | bool                | Check if label matches regex                                 |
| .LabelExistsAllOf | Function | labels ...string           | bool                | Check if all of the given labels exists                      |
| .LabelExistsAnyOf | Function | labels ...string           | bool                | Check if any of the given labels exists                      |
| ---               | ---      | ---                        | ---                 | ---                                                          |
| .Env              | Variable | none                       | map\[string\]string | A map of all environment variables and their values          |
| .EnvGet           | Function | key string                 | string              | Get a environment variable                                   |
| .EnvExists        | Function | key string                 | bool                | Check if a environment variable exists                       |
| .EnvOrDefault     | Function | key string, default string | string              | Get a environment variable or the default                    |
| .EnvEquals        | Function | key string, compare string | bool                | Check if a environment variables equals a given value        |
| .EnvMatches       | Function | key string, pattern string | bool                | Check if a environment variables equals a given value        |


### Swarm labels
If `TSD_DOCKER_TAGS_FROM_SWARM` is set to true then all of these labels are also added to the list of tags.
- com.docker.stack.namespace
- com.docker.swarm.node.id
- com.docker.swarm.service.id
- com.docker.swarm.service.name
- com.docker.swarm.task
- com.docker.swarm.task.id
- com.docker.swarm.task.name


## Raw variables
These variables contain all available data about the container itself. Shortcut methods above are much easier to use and should provide all required information for typical use cases.

| Name       | Type     | Params | Return                   | Description                  |
| ---        | ---      | ---    | ---                      | ---                          |
| .Container | Variable | none   | types.Container          | docker api container         |
| .Image     | Variable | none   | types.ImageSummary       | docker api image             |

## Examples
**Filter by image:**
```
{{- if .MatchImage "nginx" }}

# Read Nginx's basic status information (ngx_http_stub_status_module)
[[inputs.nginx]]
...

{{end}}
```

**Filter by label:**
```
{{- if .LabelExists "my.custom.nginx.detector.label" }}

# Read Nginx's basic status information (ngx_http_stub_status_module)
[[inputs.nginx]]
...

{{end}}
```

**Use variables:**
```
{{- if .MatchImage "nginx" }}

# Read Nginx's basic status information (ngx_http_stub_status_module)
[[inputs.nginx]]
  urls = ["http://{{.BridgeIP}}:{{.ConfigOrDefault "nginx_status_port" "8888" -}}
           {{- .ConfigOrDefault "nginx_status_url" "/status/nginx"}}"]

{{end}}
```

**Apply provided tags:**
```
{{- if .MatchImage "nginx" }}

# Read Nginx's basic status information (ngx_http_stub_status_module)
[[inputs.nginx]]
  urls = ["http://{{.BridgeIP}}:{{.ConfigOrDefault "nginx_status_port" "8888" -}}
           {{- .ConfigOrDefault "nginx_status_url" "/status/nginx"}}"]

  # add automatically discovered tags
  [inputs.nginx.tags]
  {{ range $key, $value := .Tags }}
  {{ $key }} = "{{ $value }}"
  {{ end }}

{{end}}
```
