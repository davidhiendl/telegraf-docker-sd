#!/bin/bash

# ensure default dirs exist
mkdir -p /etc/telegraf/telegraf.d
mkdir -p /etc/telegraf/sd-tpl.d

# write access to telegraf configuration folder is required
# - pre-existing configuration files may be read-only as well (no -R needed)
# - own configuration may be read only
chown telegraf:telegraf /etc/telegraf/telegraf.d

# ensure telegraf has access to docker socket
usermod -a -G docker telegraf

# reload init script
systemctl daemon-reload
systemctl enable telegraf-docker-sd.service

# start daemon
systemctl start telegraf-docker-sd.service
