#!/bin/bash

CFG_DIR=/etc/telegraf

# ensure default dirs exist
mkdir -p $CFG_DIR/telegraf.d
mkdir -p $CFG_DIR/sd-tpl.d

# ensure service env config file exists, this file may contain sensitive data, secure access to it
CFG_ENV_FILE=$CFG_DIR/telegraf-docker-sd.env
touch $CFG_ENV_FILE
chown telegraf:telegraf $CFG_ENV_FILE
chmod 640 $CFG_ENV_FILE


# write access to telegraf configuration folder is required
# - pre-existing configuration files may be read-only as well (no -R needed)
# - own configuration may be read only
chown telegraf:telegraf $CFG_DIR/telegraf.d

# ensure telegraf has access to docker socket
usermod -a -G docker telegraf

# reload init script
systemctl daemon-reload
systemctl enable telegraf-docker-sd.service

# start daemon
systemctl start telegraf-docker-sd.service
