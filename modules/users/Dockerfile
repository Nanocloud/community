FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

COPY ./ /go/build/users

WORKDIR /go/src/github.com/Nanocloud/community/modules/users

RUN go get ./... && go build

CMD ["./users"]
