FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> 


RUN apt-get update && \
    apt-get -y install git qemu-system-x86 sshpass

RUN mkdir -p /go/build/iaas

COPY ./ /go/src/github.com/Nanocloud/community/modules/iaas

WORKDIR /go/src/github.com/Nanocloud/community/modules/iaas

RUN mkdir -p /go/src/github.com/labstack && git clone git@github.com:labstack/echo.git -b v1 /go/src/github.com/labstack/echo
RUN go get ./... && go build

RUN mkdir -p /var/lib/nanocloud/downloads && mkdir -p /var/lib/nanocloud/images && mkdir -p /var/lib/nanocloud/pid
COPY ./scripts/ /var/lib/nanocloud/scripts/

EXPOSE 8080
CMD ["./iaas"]
