FROM golang:1.6
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> 


RUN apt-get update && \
    apt-get -y install git qemu-system-x86 sshpass genisoimage

RUN mkdir -p /go/build/iaas

COPY ./ /go/src/github.com/Nanocloud/community/modules/iaas

WORKDIR /go/src/github.com/Nanocloud/community/modules/iaas

RUN ./install.sh && go build

RUN mkdir -p /var/lib/nanocloud/downloads && mkdir -p /var/lib/nanocloud/images && mkdir -p /var/lib/nanocloud/pid && mkdir -p /var/lib/nanocloud/disk
COPY ./scripts/ /var/lib/nanocloud/scripts/
COPY ./disk/ /var/lib/nanocloud/disk/

EXPOSE 8080
EXPOSE 9090
CMD ["./iaas"]
