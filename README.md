# Telegraf Docker Service Discovery

## Description
`telegraf-docker-sd` is a companion for the
[github.com/influxdata/telegraf](https://github.com/influxdata/telegraf)
agent with the goal to support flexible configuration in order to detect
and monitor services running inside docker containers without the need
to manually configure the required inputs.

Instead of configuring every service/container the configuration is
generated via templates that are matched against containers to build
the configuration dynamically. This allows for great flexibility and
the ability to monitor docker containers that are created dynamically by
orchestration frameworks like Swarm, K8Ns, ...

## Example configuration
[Full Template Documentation](doc/TEMPLATE.md) \
By using GO Templates an enormous amount of flexibility can be achieved
when creating templates. See the full documentation for a list of
available methods and variables as well as examples.
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

## Pre-configured templates
others must be configured manually (pull requests welcome)
- MySQL
- NGINX
- PHP-FPM

## Configuration Variables
| Variable             | Default                     | Description                                                                                 |
| ---                  | ---                         | ---                                                                                         |
| TSD_TEMPLATE_DIR     | /etc/telegraf/conf.sd-tpl.d | Where configurations templates are taken from                                               |
| TSD_CONFIG_DIR       | /etc/telegraf/conf.d        | Where configurations are written to, the telegraf config directory                          |
| TSD_TAG_SWARM_LABELS | true                        | If docker swarm labels should be imported as tags. See `Container Detection > Swarm Labels` |
| TSD_TAG_LABELS       | none                        | A list of comma separated labels that should be added as tags                               |
| TSD_QUERY_INTERVAL   | 15                          | Interval in seconds between querying of the docker api for changes                          |

## Container Detection
The discovery routine automatically scans the docker host for all of it's running containers and applies each template against the container.
It is the templates responsibility to determine if the container matches the template. This allows for maximum flexibility.

### Automatically collected tags and configurations
- Any container labels starting with `telegraf.sd.tags.*` will be added as to the tag list for this container.
- Any container labels starting with `telegraf.sd.config.*` will be shortened by this prefix and added to the config map.

### Swarm labels
If `TSD_TAG_SWARM_LABELS` is set to true then all of these labels are also added to the list of tags.
- com.docker.stack.namespace
- com.docker.swarm.node.id
- com.docker.swarm.service.id
- com.docker.swarm.service.name
- com.docker.swarm.task
- com.docker.swarm.task.id
- com.docker.swarm.task.name
