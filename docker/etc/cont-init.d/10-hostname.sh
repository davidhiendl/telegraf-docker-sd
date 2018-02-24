#!/usr/bin/with-contenv sh
set -e

# find docker host hostname and expose as AGENT_HOSTNAME from /rootfs/etc
if [ -z "${AGENT_HOSTNAME+x}" ] && [ -e /rootfs/etc/hostname ]; then
    AGENT_HOSTNAME=$(cat /rootfs/etc/hostname)
    echo -n $AGENT_HOSTNAME > /var/run/s6/container_environment/AGENT_HOSTNAME
fi

