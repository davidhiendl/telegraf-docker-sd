# Telegraf Docker Service Discovery

## Description
`telegraf-docker-sd` is a lightweight (written in GO) companion for the
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
By using GO Templates an enormous amount of flexibility can be achieved
when creating templates. See the full documentations for a list of
available methods and variables as well as examples.

**[Main Template Documentation](doc/MAIN_TEMPLATE.md)** \
File: [_telegraf.goconf](sd-tpl.d/_telegraf.goconf)
```
...
# Global tags can be specified here in key="value" format.
[global_tags]
  # import all environment variables with format "GLOBAL_TAGS_$key=$value" as tags
  {{ range $key, $value := .GlobalTagsFromEnv }}
  {{ $key }} = "{{ $value }}"
{{ end }}
...
```

**[Container Template Documentation](doc/CONTAINER_TEMPLATE.md)** \
File: [nginx.goconf](sd-tpl.d/nginx.goconf)
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

## Ideas / New template methods / Issues / ... ?
Feel free to send me a PR or open an issue. I'm open for suggestions / improvements.

**Want to add custom template methods?** \
Just add new Receivers to the [sdtemplate.Params](sdtemplate/params.go) struct
```go
func (params *Params) YourCustomTemplateMethod(arg1 string, arg2 string, <<whatever>>) string {
    // do something useful
    return "somevalue"
}
```

## Pre-configured templates
others must be configured manually (pull requests welcome)
- MySQL
- NGINX
- PHP-FPM

## Docker Image
**Using pre-built image:**
```bash
docker run -ti \
    -v "$(PWD)/sd-tpl.d":/etc/telegraf/sd-tpl.d \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -v /proc:/rootfs/proc:ro \
    -v /sys:/rootfs/sys:ro \
    -v /etc:/rootfs/etc:ro \
    -e TSD_TAG_LABELS="my.custom.label.a,some.other.label.to.use.as.tags,..." \
    dhswt/telegraf-docker-sd:stable
```

**Building the image yourself:** \
The entire build (including building the binary) is included in the [Dockerfile](Dockerfile).
```bash
docker build -t yourprefix/telegraf-docker-sd:<tag>
```

## Configuration Variables
| Variable             | Default                  | Description                                                                                 |
| ---                  | ---                      | ---                                                                                         |
| TSD_TEMPLATE_DIR     | /etc/telegraf/sd-tpl.d   | Where configurations templates are taken from                                               |
| TSD_CONFIG_DIR       | /etc/telegraf/telegraf.d | Where configurations are written to, the telegraf config directory                          |
| TSD_TAG_SWARM_LABELS | true                     | If docker swarm labels should be imported as tags. See `Container Detection > Swarm Labels` |
| TSD_TAG_LABELS       | none                     | A list of comma separated labels that should be added as tags                               |
| TSD_QUERY_INTERVAL   | 15                       | Interval in seconds between querying of the docker api for changes                          |

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

## Dependencies
- GO >= 1.8 (should work with 1.7, required by docker lib, untested)
- influxdata/telegraf >= 0.10.1 (re-loading via SIGHUP is required and was implemented at that version)
- jordansissel/fpm >= 1.9.3 (debian packaging, tested only with 1.9.3, might work with earlier versions)
