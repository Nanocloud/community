FROM golang:1.6
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com>

COPY ./ /go/src/github.com/Nanocloud/community/nanocloud
WORKDIR /go/src/github.com/Nanocloud/community/nanocloud

RUN ./install.sh && go build

EXPOSE 8080
CMD ["./nanocloud"]
