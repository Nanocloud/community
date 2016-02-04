FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com>


RUN apt-get update && apt-get install sshpass
RUN mkdir -p /go/build/apps

COPY ./ /go/src/github.com/Nanocloud/community/modules/apps

WORKDIR /go/src/github.com/Nanocloud/community/modules/apps
RUN go get ./... && go build

CMD ["./apps"]
