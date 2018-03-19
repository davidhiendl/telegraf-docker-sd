# Telegraf Docker Service Discovery

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

## 0.3.0-alpha release
This release is a major refactor and rewrite with the goal to support
multiple discovery backends (kubernetes is WIP).

Most functions are still supported however the template configuration
has changed and now uses YAML files to add more structured data with an
embedded golang template.

## Example configuration
By using GO Templates an enormous amount of flexibility can be achieved
when creating templates. See the full documentations for a list of
available methods and variables as well as examples.

**[Main Template Documentation](doc/MAIN_TEMPLATE.md)** \
File: [_telegraf.yaml](sd-tpl.d/_global_telegraf.yaml)
```yaml
backend: global
template: |
    # Global tags can be specified here in key="value" format.
    [global_tags]
    # import all environment variables with format "GLOBAL_TAGS_$key=$value" as tags
    {{ as_key_value_map .Tags 2 }}

    # Configuration for telegraf agent
    [agent]
      ## Default data collection interval for all inputs
      interval = "10s"
...
```

### Backend: Docker
Monitor containers based on labels. Works for Swarm and standalone containers.

Docs: [Docker Backend Docs](doc/backend-docker/README.md) \
Example: [docker_nginx.yaml](sd-tpl.d/docker_nginx.yaml)
```yaml
backend: docker
template: |
    {{- if .MatchImage "nginx" }}

    # Read Nginx's basic status information (ngx_http_stub_status_module)
    [[inputs.nginx]]
      # An array of Nginx stub_status URI to gather stats.
      urls = ["http://{{.BridgeIP}}:{{.ConfigOrDefault "nginx_status_port" "8888" -}}
               {{- .ConfigOrDefault "nginx_status_url" "/status/nginx"}}"]

      # HTTP response timeout (default: 5s)
      response_timeout = "5s"

      # add discovered tags
      [inputs.nginx.tags]
    {{ as_key_value_map .Tags 2 }}

    {{ end -}}
```

### Backend: Kubernetes
Monitor pods based on annotations and labels. It is also possibly to use the Telegraf "prometheus" input to collect metrics from various prometheus exporters like kube-state-metrics for example:

Docs: [Kubernetes Backend Docs](doc/backend-kubernetes/README.md) \
Example: [kubernetes_kube-state-metrics.yaml](sd-tpl.d/kubernetes_kube-state-metrics.yaml)
```yaml
backend: kubernetes
template: |
    {{- if .AnnotationEquals "telegraf.sd.tags/application" "kube-state-metrics" }}

    [[inputs.prometheus]]
      ## An array of urls to scrape metrics from.
      urls = ["http://{{ .TargetIP }}:{{ .ConfigOrDefault "metrics-port" "9100" }}{{ .ConfigOrDefault "metrics-path" "/metrics" }}"]

      [inputs.prometheus.tags]
    {{ as_key_value_map .Tags 2 }}

    {{ end -}}
```

## Ideas / New template methods / Issues / ... ?
Feel free to send me a PR or open an issue. I'm open for suggestions / improvements.

**Want to add custom template methods?** \
Just add new Receivers to the tracking structs:
- [backend.docker.TrackedContainer](app/backend/docker/tracked_container.go)
- [backend.kubernetes.TrackedPod](app/backend/kubernetes/tracked_pod.go)
```go
func (td *TrackedContainer) YourCustomTemplateMethod(arg1 string, arg2 string, <<whatever>>) string {
    // do something useful
    return "somevalue"
}
```

## Pre-configured templates
others must be configured manually (pull requests welcome)
- MySQL
- NGINX
- Aerospike

## Docker Image
**Using pre-built image:**
```bash
docker run -ti \
    -v "$(PWD)/sd-tpl.d":/etc/telegraf/sd-tpl.d \
    -v /var/run/docker.sock:/var/run/docker.sock:ro \
    -v /proc:/rootfs/proc:ro \
    -v /sys:/rootfs/sys:ro \
    -v /etc:/rootfs/etc:ro \
    -v /var/run/utmp:/var/run/utmp:ro \
    -e TSD_TAG_LABELS="my.custom.label.a,some.other.label.to.use.as.tags,..." \
    dhswt/telegraf-docker-sd
```

**Building the image yourself:** \
The entire build (including building the binary) is included in the [Dockerfile](Dockerfile).
```bash
docker build -t yourprefix/telegraf-docker-sd:<tag>
```

## Configuration Variables
**TODO needs update, ref:** [ConfigSpec](app/config/config.go) [DockerConfigSpec](app/backend/docker/config.go)

| Variable             | Default                  | Description                                                                                 |
| ---                  | ---                      | ---                                                                                         |
| TSD_TEMPLATE_DIR     | /etc/telegraf/sd-tpl.d   | Where configurations templates are taken from                                               |
| TSD_CONFIG_DIR       | /etc/telegraf/telegraf.d | Where configurations are written to, the telegraf config directory                          |
| TSD_QUERY_INTERVAL   | 15                       | Interval in seconds between querying of the docker api for changes                          |

## Dependencies
- GO >= 1.9
- influxdata/telegraf >= 0.10.1 (re-loading via SIGHUP is required and was implemented at that version)
- jordansissel/fpm >= 1.9.3 (debian packaging, tested only with 1.9.3, might work with earlier versions)
