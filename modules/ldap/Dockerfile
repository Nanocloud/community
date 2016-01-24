FROM golang:1.5.2
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com>


RUN apt-get update && apt-get -y install libldap2-dev sshpass
RUN mkdir -p /go/build/ldap

COPY ./ /go/build/ldap

WORKDIR /go/build/ldap
ENV GOBIN $GOPATH/build/ldap
RUN go get ./... && go build

CMD ["./ldap"]

