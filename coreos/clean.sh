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

BUILD_ARTIFACTS="
    coreos.key
    coreos.key.pub
    coreos.qcow2
    coreos_production_qemu_image.img
    coreos_production_qemu.sh
"

for ARTIFACT in ${BUILD_ARTIFACTS}; do
    ARTIFACT_FILENAME=${CURRENT_DIR}/${ARTIFACT}
    if [ -f "${ARTIFACT_FILENAME}" ]; then
        echo "Removing ${ARTIFACT_FILENAME}"
        rm -rf "${ARTIFACT_FILENAME}"
    fi
done

NC_QEMU_PID=$(pgrep -fa qemu-system-x86 | awk '/coreos_production_qemu/ { print $1; }')
if [ -n "${NC_QEMU_PID}" ]; then
    for PID in ${NC_QEMU_PID}; do
        echo "Killing process ${PID}"
        kill "${PID}"
    done
fi
