#!/bin/bash

DOCKER_REPO=dhswt/telegraf-docker-sd
BINARY_PATH=./dist
BINARY_NAME=telegraf-docker-sd
BINARY_TARGET="${BINARY_PATH}/${BINARY_NAME}"

function build {
   set-gopath
   echo "building binary..."
   echo "> target: ${BINARY_TARGET}"
   go build -i -ldflags="-s -w" -o "${BINARY_TARGET}" ./main.go
   upx "${BINARY_TARGET}"
}

function build-dev {
   set-gopath
   echo "building binary..."
   echo "> target: ${BINARY_TARGET}"
   go build -i -ldflags="-s -w" -o "${BINARY_TARGET}" ./main.go
}

function package-deb {
    DEB_VERSION=$1

    # version is required
    if [ -z "$DEB_VERSION" ]; then
        echo "Cannot build package, missing DEB_VERSION variable"
        exit 1
    fi

    echo "building version: ${DEB_VERSION}"

    # build binary
    build

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

}

function set-gopath {
    LOCAL_GO_PATH=$(realpath $PWD/../../../..)
    GOPATH=$GOPATH:$LOCAL_GO_PATH
}

function exec-glide {
    set-gopath
    glide "${@:1}"
}

function test-run-docker {
    echo "running with backends=docker ..."
	BACKENDS=docker \
	CONFIG_DIR=./conf.test.d \
	TEMPLATE_DIR=./sd-tpl.d \
	GLOBAL_TAGS_ABC=some-value-a \
	GLOBAL_TAGS_DEF=some-value-b \
	GLOBAL_TAGS_GHI_JKL=some-value-c \
	${BINARY_TARGET}
}

function test-run-kubernetes {
    echo "running with backends=kubernetes ..."
	BACKENDS=kubernetes \
	CONFIG_DIR=./conf.test.d \
	TEMPLATE_DIR=./sd-tpl.d \
	GLOBAL_TAGS_ABC=some-value-a \
	GLOBAL_TAGS_DEF=some-value-b \
	GLOBAL_TAGS_GHI_JKL=some-value-c \
	TSD_KUBERNETES_TAG_LABELS_WHITELIST=app,name \
	TSD_KUBERNETES_NODE_NAME_OVERRIDE=k8s-prod-worker-3 \
	${BINARY_TARGET}
}

function image {
    echo "Building telegraf-docker-sd image, this might take a long time..."; \
	docker build --squash -t $DOCKER_REPO:dev .
}

function push-dev {
    echo "Pushing telegraf-docker-sd image, this might take a long time..."; \
	docker push $DOCKER_REPO:dev
}

case "$1" in
    build)
        build "${@:2}"
        ;;

    build-dev)
        build-dev "${@:2}"
        ;;

    package-deb)
        package-deb "${@:2}"
        ;;

    exec-glide)
        exec-glide "${@:2}"
        ;;

    image)
        image "${@:2}"
        ;;

    push-dev)
        push-dev "${@:2}"
        ;;

    test-run-docker)
        test-run-docker "${@:2}"
        ;;

    test-run-kubernetes)
        test-run-kubernetes "${@:2}"
        ;;

    *)
        echo $"Usage: $0 {build|build-dev|package-dev|glide|image|test-run-docker|test-run-kubernetes}"
        exit 1
esac
