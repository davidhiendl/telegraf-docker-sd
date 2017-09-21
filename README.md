# Telegraf Docker Service Discovery

## Template data struct
`TODO`  add full description of all variables & methods available to the struct (basically everything in `Params` struct

## Configuration Variables
| Variable             | Default                     | Description                                                                                 |
| ---                  | ---                         | ---                                                                                         |
| TSD_TEMPLATE_DIR     | /etc/telegraf/conf.sd-tpl.d | Where configurations templates are taken from                                               |
| TSD_CONFIG_DIR       | /etc/telegraf/conf.d        | Where configurations are written to, the telegraf config directory                          |
| TSD_TAG_SWARM_LABELS | true                        | If docker swarm labels should be imported as tags. See `Container Detection > Swarm Labels` |
| TSD_TAG_LABELS       | none                        | A list of comma separated labels that should be added as tags                               |
| TSD_QUERY_INTERVAL   | 10                          | Interval in seconds between querying of the docker api for changes                          |

## Container Detection
The discovery routine automatically scans the docker host for all of it's running containers and applies each template against the container.
It is the templates responsibility to determine if the container matches the template. This allows for maximum flexibility.

### Automatically collected tags and configurations
- Any container labels starting with `telegraf.sd.tags.*` will be added as to the tag list for this container.
- Any container labels starting with `telegraf.sd.config.*` will be shortened by this prefix and added to the config map.

### Swarm labels
If `TSD_TAG_SWARM_LABELS` is set to true then any of these labels are also added to the list of tags.
- com.docker.stack.namespace
- com.docker.swarm.node.id
- com.docker.swarm.service.id
- com.docker.swarm.service.name
- com.docker.swarm.task
- com.docker.swarm.task.id
- com.docker.swarm.task.name

## References
- [Telegraf Issue: Cannot peronalize tags on inputs per-plugin instance](https://github.com/influxdata/telegraf/issues/662)
- [GO Template Reference](https://golang.org/pkg/text/template/)

