FROM golang:1.5.2
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com>

RUN apt-get update && \
    apt-get -y install sshpass

RUN mkdir -p /go/build/nanocloud
WORKDIR /go/build/nanocloud
ENV GOBIN $GOPATH/build/nanocloud

COPY ./ /go/build/nanocloud

RUN go get ./... && go build

EXPOSE 8080
CMD ["./nanocloud"]