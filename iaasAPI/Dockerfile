FROM golang:1.5.2
MAINTAINER \
  Romain Soufflet <romain.soufflet@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \
  William Riancho <william.riancho@nanocloud.com> \


RUN apt-get update && \
    apt-get -y install git qemu-system-x86

RUN mkdir -p /go/build/iaasAPI
COPY ./ /go/build/iaasAPI
WORKDIR /go/build/iaasAPI
ENV GOBIN $GOPATH/build/iaasAPI

RUN go get ./... && go build

EXPOSE 8082
CMD ["./iaasAPI"]
