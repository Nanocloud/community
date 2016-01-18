FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \


RUN mkdir -p /go/build/users

COPY ./ /go/build/users

WORKDIR /go/build/users
ENV GOBIN $GOPATH/build/users
RUN go get ./... && go build

CMD ["./users"]
