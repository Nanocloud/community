# Plaza

Windows agent to provision and controle Nanocloud actions

# Installation

Run the following command to build Plaza:

```
./install.sh
./build.sh
```

## Installation via docker

You can also build and compile it with docker with

```
./build.sh docker
```

## Configure your build

By default, build script will generate a *plaza.exe* file for windows
platforms. If you wish to build plaza for linux, run the following command

```
export GOOS=linux
export GOARCH=amd64
./build.sh
```
