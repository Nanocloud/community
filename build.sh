#!/bin/bash -e
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

DATE_FMT="+%Y/%m/%d %H:%M:%S"


COMMAND=${1}

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
    "${CURRENT_DIR}/windows/build-windows.sh"
fi

COREOS_QCOW2_FILENAME="${CURRENT_DIR}/coreos/coreos.qcow2"
if [ -f "${COREOS_QCOW2_FILENAME}" -o "${NANOCLOUD_SKIP_COREOS}" = "true" ]; then
    echo "$(date "${DATE_FMT}") Skip CoreOS build"
else
    (
        cd coreos/
        ./build-coreos.sh
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
    [ -d iaas/scripts ] || mkdir -p iaas/scripts
    [ -d iaas/sockets ] || mkdir -p iaas/sockets

    if [ -d data ]; then
      echo "$(date "${DATE_FMT}") Erasing previous build artifact"
      rm -rf data
    fi

    (
        echo "$(date "${DATE_FMT}") Building Iaas API"
        cd iaasAPI
        go build -o "${IAAS_DIR}/scripts/api"
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
