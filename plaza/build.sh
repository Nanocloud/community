#!/bin/sh

COMMAND=${1}

if [ "${COMMAND}" = "docker" ]; then
    docker build -t nanocloud/plaza .
    docker run -d --name plaza nanocloud/plaza
    docker cp plaza:/go/src/github.com/Nanocloud/community/plaza/plaza.exe .
    docker kill plaza
    docker rm plaza
    docker rmi nanocloud/plaza
else
    export GOOS=windows
    export GOARCH=amd64

    go build
fi
