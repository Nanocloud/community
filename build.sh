#!/bin/bash
#
# Nanocloud Community, a comprehensive platform to turn any application
# into a cloud solution.
#
# Copyright (C) 2015 Nanocloud Software
#
# This file is part of Nanocloud community.
#
# Nanocloud community is free software; you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# Nanocloud community is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.


SCRIPT_FULL_PATH=$(readlink -e "${0}")
CURRENT_DIR=$(dirname "${SCRIPT_FULL_PATH}")
IAAS_DIR="${CURRENT_DIR}/iaas"

TOOLS_DIR="${IAAS_DIR}/tools/bin/"
DATE_FMT="+%Y/%m/%d %H:%M:%S"


COMMAND=${1}
shift

if [ "${COMMAND}" = "windows" ]; then
    NANOCLOUD_SKIP_WINDOWS="false"
    NANOCLOUD_SKIP_COREOS="true"
    NANOCLOUD_SKIP="true"
elif [ "${COMMAND}" = "coreos" ]; then
    NANOCLOUD_SKIP_WINDOWS="true"
    NANOCLOUD_SKIP_COREOS="false"
    NANOCLOUD_SKIP="true"
elif [ "${COMMAND}" = "nanocloud" ]; then
    NANOCLOUD_SKIP_WINDOWS="true"
    NANOCLOUD_SKIP_COREOS="true"
    NANOCLOUD_SKIP="false"
fi


WINDOWS_QCOW2_FILENAME="${CURRENT_DIR}/windows/output-windows-2012R2-qemu/windows-server-2012R2-amd64.qcow2"
if [ -f "${WINDOWS_QCOW2_FILENAME}" -o "${NANOCLOUD_SKIP_WINDOWS}" = "true" ]; then
    echo "$(date "${DATE_FMT}") Skip Windows build"
else
    (
        cd windows/
        packer build --only=windows-2012R2-qemu windows_2012_r2.json
    )
fi

COREOS_QCOW2_FILENAME="${CURRENT_DIR}/coreos/coreos.qcow2"
if [ -f "${COREOS_QCOW2_FILENAME}" -o "${NANOCLOUD_SKIP_COREOS}" = "true" ]; then
    echo "$(date "${DATE_FMT}") Skip CoreOS build"
else
    (
        cd coreos/
        source build-coreos.sh
    )
fi

NANOCLOUD_OUTPUT="${CURRENT_DIR}/nanocloud"
if [ -f "${NANOCLOUD_OUTPUT}" -o "${NANOCLOUD_SKIP}" = "true" ]; then
    echo "$(date "${DATE_FMT}") Skip Nanocloud build"
else
    echo "# Building Nanocloud"
    echo "$(date "${DATE_FMT}") Creating Iaas directories"

    [ -d iaas/images ] || mkdir -p iaas/images
    [ -d iaas/logs ] || mkdir -p iaas/logs
    [ -d iaas/pid ] || mkdir -p iaas/pid
    [ -d iaas/qemu ] || mkdir -p iaas/qemu
    [ -d iaas/scripts ] || mkdir -p iaas/scripts
    [ -d iaas/sockets ] || mkdir -p iaas/sockets
    [ -d iaas/tools ] || mkdir -p iaas/tools

    if [ -d data ]; then
      echo "$(date "${DATE_FMT}") Erasing previous build artifact"
      rm -rf data
    fi

    echo "$(date "${DATE_FMT}") Compiling Tools"

    echo -n "$(date "${DATE_FMT}") ## tunctl…  "
    (
    mkdir uml-utilities
    cd uml-utilities
    apt-get source uml-utilities
    cd "$(ls -d ./*/)"
    cd tunctl

    sed -i '3 a\LDFLAGS ?= -static' Makefile
    sed -i '/CFLAGS) -o / c\\t$(CC) $(CFLAGS) -o $(BIN) $(OBJS) $(LDFLAGS)' Makefile

    make

    cp tunctl "${TOOLS_DIR}"
    ) > /dev/null 2>&1
    echo "BUILD OK"
    rm -rf uml-utilities

    echo -n "$(date "${DATE_FMT}") ## screen…  "
    (
    mkdir screen
    cd screen
    apt-get source screen
    cd "$(ls -d ./*/)"
    ./autogen.sh
    ./configure LDFLAGS="-static"
    make

    cp screen "${TOOLS_DIR}"
    ) > /dev/null 2>&1
    echo "BUILD OK"
    rm -rf screen

    (
    echo "$(date "${DATE_FMT}") Building Iaas API"
    cd iaasAPI
    go build -a -o "${IAAS_DIR}/scripts/api"
    )

    (
    echo "$(date "${DATE_FMT}") Creating data dir"
    [ -d data ] || mkdir data

    echo "$(date "${DATE_FMT}") Compressing Nanocloud installation"
    cd iaas/
    if [ -f ../data/iaas.tar.gz ]; then
        echo "ERROR: file $(readlink -e ../data/iaas.tar.gz) already exists"
        echo "To avoid data loss, this script will end here"
        exit 2
    fi
    tar -cpzf ../data/iaas.tar.gz ./*
    )

    (
    echo "$(date "${DATE_FMT}") Bind data"
    go-bindata data/iaas.tar.gz
    mv bindata.go installer/

    cd installer
    echo -n "$(date "${DATE_FMT}") Building…"
    go build -a -o ../nanocloud
    if [ "${?}" == "0" ]; then
      echo "OK"
    else
      echo "BUILD FAILURE"
    fi
    )

    echo "Building complete, use the following command as *root* to use it"
    echo "    > ./nanocloud.sh local"
fi
