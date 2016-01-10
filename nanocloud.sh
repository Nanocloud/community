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

NANOCLOUD_DIR=${NANOCLOUD_DIR:-"${CURRENT_DIR}/installation_dir"}
DOCKER_COMPOSE_BUILD_OUTPUT="${CURRENT_DIR}/dockerfiles/build_output"

WINDOWS_QCOW2_FILENAME="${CURRENT_DIR}/windows/output-windows-2012R2-qemu/windows-server-2012R2-amd64.qcow2"
if [ -f "${WINDOWS_QCOW2_FILENAME}" ]; then
  echo "$(date "${DATE_FMT}") Local Windows image found, copying"
  cp "${WINDOWS_QCOW2_FILENAME}" "${CURRENT_DIR}/installation_dir/images/windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64.qcow2"
fi

bash $NANOCLOUD_DIR/scripts/start.sh
