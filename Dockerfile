FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com>

RUN go get -u github.com/constabulary/gb/...

COPY . /app
RUN gb build
WORKDIR /app

CMD ["./bin/apps"]
