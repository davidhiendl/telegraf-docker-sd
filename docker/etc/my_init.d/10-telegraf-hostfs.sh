#!/bin/bash
set -e

# set telegraf host variables to defaults

HOST_PROC=${HOST_PROC:-/rootfs/proc}
echo $HOST_PROC > /etc/container_environment/HOST_PROC

HOST_SYS=${HOST_SYS:-/rootfs/sys}
echo $HOST_SYS > /etc/container_environment/HOST_SYS

HOST_MOUNT_PREFIX=${HOST_SYS:-/rootfs}
echo $HOST_MOUNT_PREFIX > /etc/container_environment/HOST_MOUNT_PREFIX
