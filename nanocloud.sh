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

if [ -f "${DOCKER_COMPOSE_BUILD_OUTPUT}" ]; then
    echo "$(date "${DATE_FMT}") Starting nanocloud containers from local build"
    docker-compose --file "${CURRENT_DIR}/dockerfiles/docker-compose.yml" --x-networking up -d
else
    echo "$(date "${DATE_FMT}") Starting nanocloud containers from docker hub"
    docker-compose --x-networking up -d
fi

NANOCLOUD_STATUS=""
echo "$(date "${DATE_FMT}") Testing connectivity"
for run in {1..60} ; do
    if [ "${NANOCLOUD_STATUS}" != "200" ]; then
	CURL_CMD=$(which curl)
	WGET_CMD=$(which wget)
	if [ -n "${CURL_CMD}" ]; then
            NANOCLOUD_STATUS=$(curl --output /dev/null --insecure --silent --write-out '%{http_code}\n' "https://localhost")
	elif [ -n "${WGET_CMD}" ]; then
            NANOCLOUD_STATUS=$(LANG=C wget --no-check-certificate "https://localhost" -O /dev/null 2>&1 | awk '/^HTTP/ { print $6 ;}')
	fi
    else
	break ;
    fi
    sleep 1
done

if [ "${NANOCLOUD_STATUS}" != "200" ]; then
    echo "$(date "${DATE_FMT}") Cannot connect to Nanocloud"
    exit 1
fi

WINDOWS_QCOW2_FILENAME="${CURRENT_DIR}/windows/output-windows-2012R2-qemu/windows-server-2012R2-amd64.qcow2"
if [ -f "${WINDOWS_QCOW2_FILENAME}" ]; then
  echo "$(date "${DATE_FMT}") Local Windows image found, copying"
  cp "${WINDOWS_QCOW2_FILENAME}" "${CURRENT_DIR}/installation_dir/images/windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64.qcow2"
fi

echo "$(date "${DATE_FMT}") Setup complete"
echo "$(date "${DATE_FMT}") You can now manage your platform on : https://localhost"
echo "$(date "${DATE_FMT}") Default admin credential:"
printf "%s \tEmail: admin@nanocloud.com\n" "$(date "${DATE_FMT}")"
printf "%s \tPassword: admin\n" "$(date "${DATE_FMT}")"
echo ""
echo "$(date "${DATE_FMT}") Use the following commands to start, stop or get status information"
echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/start.sh"
echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/stop.sh"
echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/status.sh"
