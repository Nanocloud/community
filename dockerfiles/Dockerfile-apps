FROM golang:1.5.1
MAINTAINER \
  William Riancho <william.riancho@nanocloud.com>

RUN go get -u github.com/constabulary/gb/...
RUN mkdir -p /build/

RUN git clone --depth=1 --recursive https://github.com/Nanocloud/apps.git /build/apps

WORKDIR /build/apps
RUN mkdir bin && cp ./src/apps/scripts/exec.sh ./bin
RUN gb build

CMD ["./bin/apps"]
