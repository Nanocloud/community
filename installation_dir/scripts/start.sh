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

COMMUNITY_CHANNEL=$(cat ${CHANNEL_FILE})

if [ "${COMMUNITY_CHANNEL}" = "" ]; then
    COMMAND=${1}

    if [ "${COMMAND}" = "indiana" ]; then
	COMMUNITY_CHANNEL="indiana"
    elif [ "${COMMAND}" = "dev" ]; then
	COMMUNITY_CHANNEL="dev"
    else
	COMMUNITY_CHANNEL="stable"
    fi

    echo "$(date "${DATE_FMT}") Starting Nanocloud $COMMUNITY_CHANNEL for the first time"
    echo "$COMMUNITY_CHANNEL" > $CHANNEL_FILE
fi

if [ -z "$(which docker)" ]; then
  echo "$(date "${DATE_FMT}") Docker is missing, please install *docker*"
  exit 2
fi
if [ -z "$(which docker-compose)" ]; then
  echo "$(date "${DATE_FMT}") Docker-compose is missing, please install *docker-compose*"
  exit 2
fi
if [ -z "$(which curl)" -o -z "$(which wget)" ]; then
  echo "$(date "${DATE_FMT}") No download method found, please install *curl* or *wget*"
  exit 2
fi

if [ -f "${DOCKER_COMPOSE_BUILD_OUTPUT}" ]; then
    echo "$(date "${DATE_FMT}") Starting nanocloud containers from local build"
    if [ "${COMMUNITY_CHANNEL}" = "dev" ]; then
	docker-compose --file "${ROOT_DIR}/modules/docker-compose-dev.yml" --x-networking up -d
    else
	docker-compose --file "${ROOT_DIR}/modules/docker-compose-build.yml" --x-networking up -d
    fi
else
    echo "$(date "${DATE_FMT}") Starting nanocloud containers from docker hub $COMMUNITY_CHANNEL"
    if [ "${COMMUNITY_CHANNEL}" = "indiana" ]; then
	docker-compose --file "${ROOT_DIR}/docker-compose-indiana.yml" --x-networking up -d
    else
	docker-compose --file "${ROOT_DIR}/docker-compose.yml" --x-networking up -d
    fi
fi

NANOCLOUD_STATUS=""
echo "$(date "${DATE_FMT}") Testing connectivity"
for run in $(seq 60) ; do
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

echo "$(date "${DATE_FMT}") Nanocloud started"
printf "%s \tURL: https://localhost\n" "$(date "${DATE_FMT}")"
printf "%s \tEmail: admin@nanocloud.com\n" "$(date "${DATE_FMT}")"
printf "%s \tPassword: admin\n" "$(date "${DATE_FMT}")"
echo "$(date "${DATE_FMT}") This URL will only be accessible from this host."
echo ""
echo "$(date "${DATE_FMT}") Use the following commands as root to start, stop or get status information"
echo "$(date "${DATE_FMT}")     # $(readlink -e ${NANOCLOUD_DIR}/scripts/start.sh)"
echo "$(date "${DATE_FMT}")     # $(readlink -e ${NANOCLOUD_DIR}/scripts/stop.sh)"
echo "$(date "${DATE_FMT}")     # $(readlink -e ${NANOCLOUD_DIR}/scripts/status.sh)"
