#!/bin/bash
set -e

if [ -e /hostfs/etc/hostname ]; then
    AGENT_HOSTNAME=$(</hostfs/etc/hostname)
    echo $AGENT_HOSTNAME > /etc/container_environment/AGENT_HOSTNAME
fi
