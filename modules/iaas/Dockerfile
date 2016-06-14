FROM golang:1.6
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN mkdir -p /go/build/iaas
RUN mkdir -p /go/src/github.com/Nanocloud/community/nanocloud
WORKDIR /go/src/github.com/Nanocloud/community/nanocloud

RUN apt-get update && \
    apt-get -y install git qemu-system-x86 genisoimage

COPY install.sh /tmp/install.sh
RUN cd /tmp && ./install.sh
RUN cp -a /tmp/vendor /go/src/github.com/Nanocloud/community/modules/iaas

COPY ./ /go/src/github.com/Nanocloud/community/modules/iaas

RUN go build

EXPOSE 9090
CMD ["./iaas"]
