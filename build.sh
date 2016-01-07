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

DATE_FMT="+%Y/%m/%d %H:%M:%S"


COMMAND=${1}

if [ "${COMMAND}" = "windows" ]; then
    NANOCLOUD_SKIP_WINDOWS="false"
    NANOCLOUD_SKIP="true"
elif [ "${COMMAND}" = "nanocloud" ]; then
    NANOCLOUD_SKIP_WINDOWS="true"
    NANOCLOUD_SKIP="false"
fi

WINDOWS_QCOW2_FILENAME="${CURRENT_DIR}/windows/output-windows-2012R2-qemu/windows-server-2012R2-amd64.qcow2"
if [ -f "${WINDOWS_QCOW2_FILENAME}" -o "${NANOCLOUD_SKIP_WINDOWS}" = "true" ]; then
    echo "$(date "${DATE_FMT}") Skip Windows build"
else
    "${CURRENT_DIR}/windows/build-windows.sh"
fi

NANOCLOUD_OUTPUT="${CURRENT_DIR}/dockerfiles/build_output"
if [ -f "${NANOCLOUD_OUTPUT}" -o "${NANOCLOUD_SKIP}" = "true" ]; then
    echo "$(date "${DATE_FMT}") Skip Nanocloud build"
else
    echo "# Building Nanocloud"
    DOCKER_COMPOSE=$(which docker-compose)
    if [ -z "${DOCKER_COMPOSE}" ]; then
        echo "You need *docker-compose* to run this script, exiting"
        exit 1
    fi

    (
        cd dockerfiles
        ${DOCKER_COMPOSE} build
        if [ ${?} = 0 ]; then
            echo "0" > build_output
        fi
    )

    echo "Build completed, use the following command to use it"
    echo "    > ./nanocloud.sh"
fi
