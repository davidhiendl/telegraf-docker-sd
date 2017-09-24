#!/bin/bash
set -e

if [ -e /rootfs/etc/hostname ]; then
    AGENT_HOSTNAME=$(</rootfs/etc/hostname)
    echo $AGENT_HOSTNAME > /etc/container_environment/AGENT_HOSTNAME
fi
