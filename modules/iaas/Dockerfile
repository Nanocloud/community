FROM golang:1.6
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> 


RUN apt-get update && \
    apt-get -y install git qemu-system-x86 genisoimage

RUN mkdir -p /go/build/iaas

COPY ./ /go/src/github.com/Nanocloud/community/modules/iaas

WORKDIR /go/src/github.com/Nanocloud/community/modules/iaas

RUN ./install.sh && go build

EXPOSE 8080
EXPOSE 9090
CMD ["./iaas"]
