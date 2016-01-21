FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> 


RUN apt-get update && \
    apt-get -y install git qemu-system-x86 sshpass

RUN mkdir -p /go/build/iaas

COPY ./ /go/build/iaas

WORKDIR /go/build/iaas
ENV GOBIN $GOPATH/build/iaas
RUN go get ./... && go build

CMD ["./iaas"]
