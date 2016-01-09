FROM golang:1.5.2
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com>

RUN apt-get update && \
    apt-get -y install git qemu-system-x86

RUN go get "github.com/gorilla/rpc" && \
    go get "github.com/gorilla/rpc/json" && \
    go get "gopkg.in/yaml.v2"

RUN mkdir -p /opt/build
WORKDIR /opt/build/

RUN git clone --depth=1 --recursive https://github.com/Nanocloud/community.git

WORKDIR /opt/build/community/iaasAPI

RUN go build

EXPOSE 8082
CMD ["./iaasAPI"]
