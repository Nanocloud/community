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


NANOCLOUD_URL=${NANOCLOUD_URL:-"http://releases.nanocloud.com:8080/indiana/"}
DATE_FMT="+%Y/%m/%d %H:%M:%S"


echo "Nanocloud one-liner installer"
if [ -z "${NANOCLOUD_DIR}" ]; then
    NANOCLOUD_DIR=$(mktemp --directory /tmp/nanocloud-XXXX)
    echo "$(date "${DATE_FMT}") No installation dir specified. Creating one in ${NANOCLOUD_DIR}"
fi


download() {
    CURL_CMD=$(which curl)
    WGET_CMD=$(which wget)

    URL=${1}
    if [ -n "${CURL_CMD}" ]; then
        curl --progress-bar "${URL}"
    elif [ -n "${WGET_CMD}" ]; then
        wget --quiet "${URL}" -O -
    else
        echo "You need *curl* or *wget* to run this script, exiting"
        exit 2
    fi
}

(
    cd "${NANOCLOUD_DIR}"
    echo "$(date "${DATE_FMT}") Creating direcotories"
    mkdir -p dockerfiles/nginx/conf/certificates
    mkdir -p installation_dir/scripts
    mkdir -p installation_dir/scripts
    echo "$(date "${DATE_FMT}") Downloading artifacts"
    download "${NANOCLOUD_URL}/nanocloud.sh" > nanocloud.sh &
    download "${NANOCLOUD_URL}/docker-compose.yml" > docker-compose.yml &
    download "${NANOCLOUD_URL}/dockerfiles/docker-compose.yml" > dockerfiles/docker-compose.yml &
    download "${NANOCLOUD_URL}/dockerfiles/nginx/conf/certificates/nginx.crt" > dockerfiles/nginx/conf/certificates/nginx.crt &
    download "${NANOCLOUD_URL}/dockerfiles/nginx/conf/certificates/nginx.key" > dockerfiles/nginx/conf/certificates/nginx.key &
    download "${NANOCLOUD_URL}/dockerfiles/nginx/conf/nginx.conf" > dockerfiles/nginx/conf/nginx.conf &
    download "${NANOCLOUD_URL}/scripts/start.sh" > installation_dir/scripts/start.sh &

    wait

    echo "$(date "${DATE_FMT}") Installing…"
    chmod +x nanocloud.sh
    ./nanocloud.sh
)
