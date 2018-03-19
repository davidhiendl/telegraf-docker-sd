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
    build
    ./build/scripts.d/package-deb.sh
}

function set-gopath {
    LOCAL_GO_PATH=$(realpath $PWD/../../../..)
    GOPATH=$GOPATH:$LOCAL_GO_PATH
}

function glide {
    set-gopath
    glide "${@:2}"
}

function test-run-docker {
	BACKENDS=docker \
	CONFIG_DIR=./conf.test.d \
	TEMPLATE_DIR=./sd-tpl.d \
	GLOBAL_TAGS_ABC=some-value-a \
	GLOBAL_TAGS_DEF=some-value-b \
	GLOBAL_TAGS_GHI_JKL=some-value-c \
	${BINARY_TARGET}
}

function test-run-kubernetes {
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
        build
        ;;

    build-dev)
        build-dev
        ;;

    package-deb)
        package-deb
        ;;

    glide)
        glide
        ;;

    image)
        image
        ;;

    push-dev)
        push-dev
        ;;

    test-run-docker)
        test-run-docker
        ;;

    test-run-kubernetes)
        test-run-kubernetes
        ;;

    *)
        echo $"Usage: $0 {build|build-dev|package-dev|glide|image|test-run-docker|test-run-kubernetes}"
        exit 1
esac
