#!/bin/bash
set -e

HOST_PROC=${HOST_PROC:-/hostfs/proc}
echo $HOST_PROC > /etc/container_environment/HOST_PROC

HOST_SYS=${HOST_SYS:-/hostfs/sys}
echo $HOST_SYS > /etc/container_environment/HOST_SYS
