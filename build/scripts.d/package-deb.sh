#!/bin/bash

set -x

echo $DEB_VERSION

# version is required
if [ -z "$DEB_VERSION" ]; then
    echo "Cannot build package, missing DEB_VERSION variable"
    exit 1
fi

# prepare paths
SRC_DIR="./build/deb.d"
DST_DIR="./dist/"
BIN_SRC="${DST_DIR}/telegraf-docker-sd"
DST_PKG="${DST_DIR}/telegraf-docker-sd_${DEB_VERSION}.deb"

# check if binary exists
if [ ! -e "${BIN_SRC}" ]; then
    echo "Cannot build package, missing binary"
    exit 1
fi

# assemble package
rm -f "${DST_PKG}"
fpm -s dir \
    -t deb \
    -n telegraf-docker-sd \
    --license "MIT" \
    --maintainer "David Hiendl<david.hiendl@dhswt.de>" \
    --vendor "DHSWT" \
    --url "https://github.com/davidhiendl/telegraf-docker-sd" \
    --description "Automatic Docker service discovery for Telegraf" \
    -d "telegraf >= 0.10.1" \
    -v "${DEB_VERSION}" \
    -p "${DST_PKG}" \
    --after-install "${SRC_DIR}/after-install.sh" \
    "${SRC_DIR}/telegraf-docker-sd.service"=/etc/systemd/system/telegraf-docker-sd.service \
    "${BIN_SRC}"=/usr/local/bin/telegraf-docker-sd
