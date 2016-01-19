#!/bin/sh
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

ROOT_DIR=${CURRENT_DIR}/../..
NANOCLOUD_DIR=${NANOCLOUD_DIR:-"${ROOT_DIR}/installation_dir"}
DOCKER_COMPOSE_BUILD_OUTPUT="${ROOT_DIR}/modules/build_output"
CHANNEL_FILE=${NANOCLOUD_DIR}/channel

COMMAND=${1}

COMMUNITY_CHANNEL=$(cat ${CHANNEL_FILE})

if [ -z "$(which docker || true)" ]; then
  echo "$(date "${DATE_FMT}") Docker is missing, please install *docker*"
  exit 2
fi
if [ -z "$(which docker-compose || true)" ]; then
  echo "$(date "${DATE_FMT}") Docker-compose is missing, please install *docker-compose*"
  exit 2
fi

if [ -f "${DOCKER_COMPOSE_BUILD_OUTPUT}" ]; then
    echo "$(date "${DATE_FMT}") Stopping nanocloud containers from local build"
    docker-compose --file "${ROOT_DIR}/modules/docker-compose.yml" stop
else
    echo "$(date "${DATE_FMT}") Stopping nanocloud containers from docker hub $COMMUNITY_CHANNEL"
    if [ "${COMMUNITY_CHANNEL}" = "indiana" ]; then
	docker-compose --file "${ROOT_DIR}/docker-compose-indiana.yml" stop
    else
	docker-compose --file "${ROOT_DIR}/docker-compose.yml" stop
    fi
fi

rm -f ${NANOCLOUD_DIR}/pid/windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64.pid

echo "$(date "${DATE_FMT}") Nanocloud stopped"
echo "$(date "${DATE_FMT}") To start again Nanocloud, use :"
echo "$(date "${DATE_FMT}")     # $(readlink -e ${NANOCLOUD_DIR}/scripts/start.sh)"
