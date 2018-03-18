# Backend: Docker

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
| .Config           | Variable | none                       | map\[string\]string | A map of all configurations values                           |
| .ConfigOrDefault  | Function | key string, default string | string              | Retrieve a config value from combined sources or use default |
| .ConfigGet        | Function | key string                 | string              | Retrieve a config value from combined sources or use default |
| .MatchImageRegex  | Function | key string                 | bool                | Check container image against provided regex                 |
| .MatchImage       | Function | pattern string             | bool                | Check container image against provided literal string        |
| .LabelExists      | Function | label string               | bool                | Check container labels against provided label                |
| .LabelExistsAllOf | Function | labels ...string           | bool                | Check if container labels contain all provided labels        |
| .LabelExistsAnyOf | Function | labels ...string           | bool                | Check if container labels contain any provided labels        |
| .LabelEquals      | Function | label string, value string | bool                | Check if container labels contain any provided labels        |
| .BridgeIP         | Function | none                       | string              | The actual container bridge ip                               |
| .Tags             | Variable | none                       | map\[string\]string | A map of all computed tags and values                        |
| .Labels           | Function | none                       | map\[string\]string | A map of all labels and their values                         |
| .EnvOrDefault     | Function | label string, value string | string              | Get a environment variable or the default                    |
| .EnvGet           | Function | label string               | string              | Get a environment variable                                   |
| .EnvHas           | Function | label string               | bool                | Check if a environment variable exists                       |
| .EnvEquals        | Function | label string, value string | bool                | Check if a environment variables equals a given value        |


## Raw variables
These variables contain all available data about the container itself. Shortcut methods above are much easier to use and should provide all required information for typical use cases.

| Name       | Type     | Params | Return                   | Description                  |
| ---        | ---      | ---    | ---                      | ---                          |
| .Container | Variable | none   | types.Container          | docker api container         |
| .Bridge    | Variable | none   | network.EndpointSettings | docker api endpoint settings |
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
