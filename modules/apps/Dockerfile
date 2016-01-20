FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> 


RUN apt-get update && apt-get install sshpass
RUN mkdir -p /go/build/apps

COPY ./ /go/build/apps

WORKDIR /go/build/apps
ENV GOBIN $GOPATH/build/apps
RUN go get ./... && go build

CMD ["./apps"]
