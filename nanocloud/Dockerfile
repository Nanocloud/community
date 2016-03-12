FROM golang:1.6
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com>

RUN apt-get update && \
    apt-get -y install sshpass

COPY ./ /go/src/github.com/Nanocloud/community/nanocloud
WORKDIR /go/src/github.com/Nanocloud/community/nanocloud

RUN mkdir -p /go/src/github.com/labstack && git clone https://github.com/labstack/echo.git -b v1 /go/src/github.com/labstack/echo
RUN go get && go build

EXPOSE 8080
CMD ["./nanocloud"]
