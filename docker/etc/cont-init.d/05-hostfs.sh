#!/usr/bin/with-contenv sh
set -e

# set telegraf "HOST_" variables if present at default location
if [ -d /rootfs/proc ]; then
    HOST_PROC=${HOST_PROC:-/rootfs/proc}
    echo -n $HOST_PROC > /var/run/s6/container_environment/HOST_PROC
fi

if [ -d /rootfs/sys ]; then
    HOST_SYS=${HOST_SYS:-/rootfs/sys}
    echo -n $HOST_SYS > /var/run/s6/container_environment/HOST_SYS
fi

if [ -d /rootfs/etc ]; then
    HOST_ETC=${HOST_ETC:-/rootfs/etc}
    echo -n $HOST_ETC > /var/run/s6/container_environment/HOST_ETC
fi

