# Global templates

All templates with `backend: global` will be processed exactly once at
startup allowing for an initial, target independent configuration.

## Example

```yaml
backend: global
template: |
    [global_tags]
    {{ as_key_value_map .Tags 2 }}

    ## Agent config
    [agent]
      # Default data collection interval for all inputs
      interval = "10s"
      round_interval = true
      metric_batch_size = 1000
      metric_buffer_limit = 10000
      collection_jitter = "1s"
      flush_interval = "10s"
      flush_jitter = "5s"
      precision = ""

      # Logging configuration:
      debug = false
      quiet = false
      logfile = ""

      hostname = "{{ .EnvGet "AGENT_HOSTNAME" }}"

    ## Influxdb Output
    {{ if .EnvHas "OUTPUT_INFLUXDB_URL" }}
    [[outputs.influxdb]]
      urls = ["{{ .EnvGet "OUTPUT_INFLUXDB_URL" }}"]
      database = "{{ .EnvOrDefault "OUTPUT_INFLUXDB_DB" "telegraf" }}"
      retention_policy = ""
      write_consistency = "any"
      timeout = "5s"

      {{ if .EnvHas "OUTPUT_INFLUXDB_USER" }}
      username = "{{ .EnvGet "OUTPUT_INFLUXDB_USER" }}"
      {{ end }}

      {{ if .EnvHas "OUTPUT_INFLUXDB_PASS" }}
      password = "{{ .EnvGet "OUTPUT_INFLUXDB_PASS" }}"
      {{ end }}

    {{ end }}
```


## Available functions and variables

| Name          | Type     | Params                     | Return              | Description                                           |
| ---           | ---      | ---                        | ---                 | ---                                                   |
| .Tags         | Variable | none                       | map\[string\]string | A map of tags from environment                        |
| .EnvOrDefault | Function | label string, value string | string              | Get a environment variable or the default             |
| .EnvGet       | Function | label string               | string              | Get a environment variable                            |
| .EnvHas       | Function | label string               | bool                | Check if a environment variable exists                |
| .EnvEquals    | Function | label string, value string | bool                | Check if a environment variables equals a given value |
