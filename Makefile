.DEFAULT_GOAL := list

list:
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$' | xargs

image-repo  = dhswt/telegraf-docker-sd
source-path = /go/src/github.com/davidhiendl/telegraf-docker-sd
binary-path = ./dist
binary-name = telegraf-docker-sd

# build compressed binary using local go
binary:
	GOPATH=$$GOPATH:$$PWD/../../../../ \
	&& echo $$GOPATH \
	&& go get . \
	&& go build -i -ldflags="-s -w" -o $(binary-path)/$(binary-name) ./main.go \
	&& upx $(binary-path)/$(binary-name)

# build using local go
binary-dev:
	GOPATH=$$GOPATH:$$PWD/../../../../ \
	&& echo $$GOPATH \
	&& go get . \
	&& go build -i -o $(binary-path)/$(binary-name) ./main.go

# build image
image:
	echo "Building telegraf-docker-sd image, this might take a long time..."; \
	docker build --squash -t $(image-repo):master .

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
