# build binary first
FROM    golang:1.9.4-alpine3.7

WORKDIR /go/src/github.com/davidhiendl/telegraf-docker-sd
ENV     TELEGRAF_VERSION 1.5.2


# install upx to compress binary
RUN     apk add --no-cache \
            git \
            mercurial \
            upx \
            curl \

&&      curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# get telegraf
RUN     curl -o telegraf.tar.gz https://dl.influxdata.com/telegraf/releases/telegraf-${TELEGRAF_VERSION}-static_linux_amd64.tar.gz \
&&      tar xf telegraf.tar.gz \
&&      mv ./telegraf/telegraf / \
&&      upx /telegraf

# add sources (add manually instead of entire dir to speed up build time for non-code changes)
ADD     main.go /go/src/github.com/davidhiendl/telegraf-docker-sd/
ADD     app     /go/src/github.com/davidhiendl/telegraf-docker-sd/app
ADD     vendor  /go/src/github.com/davidhiendl/telegraf-docker-sd/vendor

# fetch remaining dependencies and build package
RUN     dep ensure \
&&      go build -i \
            -o /telegraf-docker-sd \
            -ldflags="-s -w" \
            ./main.go \

# compress binary
&&      upx /telegraf-docker-sd

# build service container next
FROM    alpine
LABEL   maintainer="David Hiendl <david.hiendl@dhswt.de>"

# add required packages
RUN     apk add --no-cache \
            iputils \
            ca-certificates \
            net-snmp-tools \
            procps \
            curl \
&&      update-ca-certificates \

# install overlay
&&      curl -L -s https://github.com/just-containers/s6-overlay/releases/download/v1.21.2.2/s6-overlay-amd64.tar.gz \
        | tar xvzf - -C / \
&&      apk del --no-cache curl

# configuration
ADD     docker/etc   /etc
ADD     sd-tpl.d     /etc/telegraf/sd-tpl.d

# clean telegraf conf
RUN     mkdir -p /etc/telegraf/telegraf.d \
&&      echo "# configured via template _telegraf.goconf\n" > /etc/telegraf/telegraf.conf \

# prepare permissions
&&      addgroup -S telegraf \
&&      adduser -s /bin/false -S -D -H telegraf \
&&      chown root:root -R /etc/telegraf \
&&      chmod -R u=rwX,g=rwX,o=rX /etc/telegraf \
&&      chown telegraf:telegraf /etc/telegraf/telegraf.d

# add binary from previous stage
COPY    --from=0 /telegraf /usr/local/bin
COPY    --from=0 /telegraf-docker-sd /usr/local/bin

ENTRYPOINT ["/init"]
