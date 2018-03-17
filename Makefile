.DEFAULT_GOAL := list

image-repo  = dhswt/telegraf-docker-sd
source-path = /go/src/github.com/davidhiendl/telegraf-docker-sd
binary-path = ./dist
binary-name = telegraf-docker-sd
go-path-local = $(shell realpath $$PWD/../../../../)

# build compressed binary using local go
dep:
	GOPATH=$$GOPATH:$(go-path-local) \
    && echo GOPATH = $$GOPATH \
    && dep ${ARGS} \

binary:
	GOPATH=$$GOPATH:$(go-path-local) \
	&& go build -i -ldflags="-s -w" -o $(binary-path)/$(binary-name) ./main.go \
	&& upx $(binary-path)/$(binary-name)

# build using local go
binary-dev:
	GOPATH=$$GOPATH:$(go-path-local) \
    && go build -i -ldflags="-s -w" -o $(binary-path)/$(binary-name) ./main.go \
	&& go build -i -o $(binary-path)/$(binary-name) ./main.go

# build image
image:
	echo "Building telegraf-docker-sd image, this might take a long time..."; \
	docker build --squash -t $(image-repo):master .

test-run:
	BACKENDS=docker \
	CONFIG_DIR=./conf.test.d \
	TEMPLATE_DIR=./sd-tpl.d \
	GLOBAL_TAGS_ABC=some-value-a \
	GLOBAL_TAGS_DEF=some-value-b \
	GLOBAL_TAGS_GHI_JKL=some-value-c \
	$(binary-path)/$(binary-name)

compose-up:
	docker-compose up

# create deb package
package-deb: binary
	./build/scripts.d/package-deb.sh

tag-push-testing:
	docker tag $(image-repo):master $(image-repo):testing && \
	docker push $(image-repo):testing

show-images:
	docker images | grep "$(image-repo)"

# Remove dangling images
clean-images:
	docker images -a -q \
		--filter "reference=$(image-repo)" \
		--filter "dangling=true" \
	| xargs docker rmi

# Remove all images
clear-images:
	docker images -a -q \
		--filter "reference=$(image-repo)" \
	| xargs docker rmi
