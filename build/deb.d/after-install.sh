#!/bin/bash

# ensure default dirs exist
mkdir -p /etc/telegraf/telegraf.d
mkdir -p /etc/telegraf/sd-tpl.d

# ensure telegraf has access to docker socket
usermod -a -G docker telegraf

# reload init script
systemctl daemon-reload

# start daemon
systemctl start telegraf-docker-sd.service
