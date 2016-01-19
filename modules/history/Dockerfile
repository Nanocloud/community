FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com> \
  Olivier Berthonneau <olivier.berthonneau@nanocloud.com> \


RUN mkdir -p /go/build/history

COPY ./ /go/build/history

WORKDIR /go/build/history
ENV GOBIN $GOPATH/build/history
RUN go get ./... && go build

CMD ["./history"]
