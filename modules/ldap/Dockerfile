FROM golang:1.5.2
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com>

RUN mkdir -p /go/build/ldap

COPY ./ /go/src/github.com/Nanocloud/community/modules/ldap

WORKDIR /go/src/github.com/Nanocloud/community/modules/ldap
RUN go get ./... && go build

CMD ["./ldap"]

