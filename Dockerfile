FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com>

RUN go get -u github.com/constabulary/gb/...

COPY . /app
WORKDIR /app
RUN gb build

CMD ["./bin/apps"]
