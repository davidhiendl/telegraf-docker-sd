# Backend: Kubernetes

The discovery routine requests information from the kubernetes api for
all resources running on the same host as itself and applies each
template against the container. It is the templates responsibility to
determine if the container matches the template. This allows for maximum
flexibility.

At the moment only pods and pod annotations can be tracked directly.


### Automatically collected tags and configurations
- Any pod annotations starting with `telegraf.sd.tags/<any-key>: <any-value>` will be shortened by this prefix and as tags for this pod.
- Any pod annotations starting with `telegraf.sd.config/<any-key>: <any-value>` will be shortened by this prefix and as config for this pod.


## Environment Configuration Parameters

| Name                                | Type   | Default     | Description                                                                                                                                                                |
| ---                                 | ---    | ---         | ---                                                                                                                                                                        |
| TSD_KUBERNETES_NODE_NAME_OVERRIDE   | bool   | true        | Overwrite the node name the agent is looking for and monitors pods on                                                                                                      |
| TSD_KUBERNETES_TAG_NAMESPACE        | bool   | true        | Include k8s_namespace=<Pod.Namespace> as tag                                                                                                                               |
| TSD_KUBERNETES_TAG_POD              | bool   | true        | Include k8s_pod=<Pod.Name> as tag                                                                                                                                          |
| TSD_KUBERNETES_AUTO_CONF_PREFIX     | string | kubernetes_ | Backend specific prefix for generated files                                                                                                                                |
| TSD_KUBERNETES_TAG_LABELS_WHITELIST | string |             | An explicit comma-separated list of labels to include as tags                                                                                                              |
| TSD_KUBERNETES_TAG_LABELS_BLACKLIST | string |             | A comma-separated list of labels to exclude from tags, if this value is non-empty all labels except the listed labels will be included. This may cause cardinality issues. |


## Available functions and variables


| Name                 | Type     | Params                     | Return              | Description                                                  |
| ---                  | ---      | ---                        | ---                 | ---                                                          |
| .TargetIP            | Function | none                       | string              | The actual container bridge ip                               |
| .Tags                | Variable | none                       | map\[string\]string | A map of all computed tags and values                        |
| ---                  | ---      | ---                        | ---                 | ---                                                          |
| .Config              | Variable | none                       | map\[string\]string | A map of all configurations values                           |
| .ConfigGet           | Function | key string                 | string              | Retrieve a config value from combined sources or use default |
| .ConfigOrDefault     | Function | key string, default string | string              | Retrieve a config value from combined sources or use default |
| .ConfigExists        | Function | key string                 | bool                | Check if a config key exists                                 |
| .ConfigEquals        | Function | key string, default string | bool                | Check if a config value equals an other value                |
| .ConfigMatches       | Function | key string, pattern string | bool                | Check if a config value matches a regex pattern              |
| ---                  | ---      | ---                        | ---                 | ---                                                          |
| .Labels              | Variable | none                       | map\[string\]string | A map of all labels and their values                         |
| .LabelGet            | Function | key string                 | bool                | Get label value                                              |
| .LabelOrDefault      | Function | key string, default string | bool                | Get label value or default                                   |
| .LabelExists         | Function | key string                 | bool                | Check if label exists                                        |
| .LabelEquals         | Function | key string, compare string | bool                | Check if label equals value                                  |
| .LabelMatches        | Function | key string, pattern string | bool                | Check if label matches regex                                 |
| .LabelExistsAllOf    | Function | labels ...string           | bool                | Check if all of the given labels exists                      |
| .LabelExistsAnyOf    | Function | labels ...string           | bool                | Check if any of the given labels exists                      |
| ---                  | ---      | ---                        | ---                 | ---                                                          |
| .Annotations         | Variable | none                       | map\[string\]string | A map of all annotations and their values                    |
| .AnnotationGet       | Function | key string                 | bool                | Get annotation value                                         |
| .AnnotationOrDefault | Function | key string, default string | bool                | Get label annotation or default                              |
| .AnnotationExists    | Function | key string                 | bool                | Check if annotation exists                                   |
| .AnnotationEquals    | Function | key string, compare string | bool                | Check if annotation equals value                             |
| .AnnotationMatches   | Function | key string, pattern string | bool                | Check if annotation matches regex                            |
| ---                  | ---      | ---                        | ---                 | ---                                                          |
| .Env                 | Variable | none                       | map\[string\]string | A map of all environment variables and their values          |
| .EnvGet              | Function | key string                 | string              | Get a environment variable                                   |
| .EnvExists           | Function | key string                 | bool                | Check if a environment variable exists                       |
| .EnvOrDefault        | Function | key string, default string | string              | Get a environment variable or the default                    |
| .EnvEquals           | Function | key string, compare string | bool                | Check if a environment variables equals a given value        |
| .EnvMatches          | Function | key string, pattern string | bool                | Check if a environment variables equals a given value        |


## Raw variables
These variables contain all available data about the container itself. Shortcut methods above are much easier to use and should provide all required information for typical use cases.

| Name  | Type     | Params | Return                 | Description        |
| ---   | ---      | ---    | ---                    | ---                |
| .Pod  | Variable | none   | k8s.io/api/core/v1.Pod | kubernetes api pod |
| .Name | Variable | none   | string                 | pod name           |
| .UID  | Variable | none   | UID (string)           | pod UID            |


## Examples
TODO examples


## TODOs
- Docs/Examples
- Tracking Services
    - directly
    - or all endspoints/targets? or both?
