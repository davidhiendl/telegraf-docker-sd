# Templates

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

{{end}}
```
