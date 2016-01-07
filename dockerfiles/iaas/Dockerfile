FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com>

RUN go get -u github.com/constabulary/gb/...
RUN mkdir -p /build/

RUN git clone --depth=1 --recursive https://github.com/Nanocloud/iaas.git /build/iaas

WORKDIR /build/iaas
RUN gb build

CMD ["./bin/iaas"]
