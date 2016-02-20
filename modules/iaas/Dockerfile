FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> 


RUN apt-get update && \
    apt-get -y install git qemu-system-x86 sshpass

RUN mkdir -p /go/build/iaas

COPY ./ /go/src/github.com/Nanocloud/community/modules/iaas

WORKDIR /go/src/github.com/Nanocloud/community/modules/iaas

RUN go get ./... && go build

EXPOSE 8080
CMD ["./iaas"]
