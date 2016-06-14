FROM golang:1.6
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com>

RUN mkdir -p /go/src/github.com/Nanocloud/community/nanocloud
WORKDIR /go/src/github.com/Nanocloud/community/nanocloud

COPY install.sh /tmp/install.sh
RUN cd /tmp && ./install.sh
RUN cp -a /tmp/vendor //go/src/github.com/Nanocloud/community/nanocloud/

COPY ./ /go/src/github.com/Nanocloud/community/nanocloud

RUN go build

EXPOSE 8080
CMD ["./nanocloud"]
