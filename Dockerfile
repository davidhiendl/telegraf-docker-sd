# build binary first
FROM	golang:1.8

WORKDIR	/go/src/github.com/davidhiendl/telegraf-docker-sd

# install upx to compress binary
RUN		apt-get update \
&&		apt-get install -y upx

# add sources (add manually instead of entire dir to speed up build time for non-code changes)
ADD		main.go		/go/src/github.com/davidhiendl/telegraf-docker-sd/
ADD		app			/go/src/github.com/davidhiendl/telegraf-docker-sd/app

# fetch remaining dependencies and build package
RUN		go get . \
&&		go build -i \
			-o /telegraf-docker-sd \
			-ldflags="-s -w" \
			./main.go \

# compress binary
&&		upx /telegraf-docker-sd

# build container next
FROM    phusion/baseimage:0.9.22
LABEL 	maintainer="David Hiendl <david.hiendl@dhswt.de>"

# install telegraf
RUN		curl -o telegraf.deb https://dl.influxdata.com/telegraf/releases/telegraf_1.4.0-1_amd64.deb \
&&		dpkg -i telegraf.deb \
&&		rm telegraf.deb

# configure services and startup
ADD		docker/services		/etc/service
ADD		docker/my_init.d	/etc/my_init.d

# update permissions
RUN		chmod 555 \
			/etc/my_init.d/* \
			/etc/service/telegraf/* \
			/etc/service/telegraf-docker-sd/*

ADD 	sd-tpl.d	/etc/telegraf/sd-tpl.d
# add binary from previous stage
COPY    --from=0 /telegraf-docker-sd /usr/local/bin
