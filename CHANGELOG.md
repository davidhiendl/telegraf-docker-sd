# Changelog

## [Unreleased]

## [0.5.0-alpha] - 2018-07-11
Improvements to handling external problems, mainly:
- Track POD IP changes in Kubernetes backend
- No longer panic/exit when the connection to Backends fails (e.g. docker service stopped, kubernetes API unreachable, ...)

### Changed
- updated kubernetes pod tracking to also detect changes in PodIP (as they are not guranteed to be constant)
- change behaviour when faced with issues during docker/k8s information gathering: no longer panic/exit directly but instead log errors and retry on next run. improves resilience to failures due to service/node reboots, network issues, ...

## [0.4.0-alpha] - 2018-04-13
Improvements to build, logging and packaging

### Changed
- fixed telegraf reloader not working correctly (regression)
- made kubernetes config file path configurable
- replaced custom logger with logrus
- improved debian packaging and added additional configuration file for telegraf-docker-sd with systemd
- improved log format
- fixed build
- switched to logrus from custom logger


## [0.3.0-alpha] - 2018-03-19
This release is a major refactor and rewrite with the goal to support
multiple discovery backends (kubernetes is WIP).

Most functions are still supported however the template configuration
has changed and now uses YAML files to add more structured data with an
embedded golang template.

### Added
- interface for handling multiple backends
- new template files based on YAML
- added multi-backend support
- added kubernetes support
- improved templating

### Changed
- refactored docker backend to support multiple backends
- switched to glide for dependencies as k8s client does not yet support dep
- reverted back to old docker go sdk as there were dependency incompatibilities with the new one which is still WIP
- fixed docker version mismatch errors by allowing the client to negotiate the api version
- massive update to documentation

### Removed
- old template files support
