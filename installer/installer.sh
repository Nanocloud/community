#!/bin/sh -e
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

SCRIPT_UID=$(id -u)

DATE_FMT="+%Y/%m/%d %H:%M:%S"
COMMUNITY_TAG="0.7.0rc1"
COMMAND=${1}

if [ -z "$(which docker || true)" ]; then
  echo "$(date "${DATE_FMT}") Docker is missing, please install *docker*"
  exit 2
fi

check_docker_compose_file () {
    DOCKERCOMPOSEFILE=${PWD}/nanocloud/docker-compose.yml

    if [ -z "$(ls -l $DOCKERCOMPOSEFILE 2> /dev/null || true)" ]; then
        echo "$(date "${DATE_FMT}") Could not find $DOCKERCOMPOSEFILE"
        exit 2
    fi
}

case $COMMAND in
    "start")
        check_docker_compose_file
        docker-compose -f ${PWD}/nanocloud/docker-compose.yml up -d
        ;;
    "stop")
        check_docker_compose_file
        docker-compose -f ${PWD}/nanocloud/docker-compose.yml kill
        ;;
    "status")
        NANOCLOUD_STATUS=$(curl --output /dev/null --insecure --silent --write-out '%{http_code}\n' "https://$(docker exec proxy hostname -I | awk '{print $1}')")
        if [ "${NANOCLOUD_STATUS}" != "200" ]; then
            echo "$(date "${DATE_FMT}") Nanocloud is *NOT* running"
        else
            echo "$(date "${DATE_FMT}") Nanocloud is running"
            printf "%s \tURL: https://localhost\n" "$(date "${DATE_FMT}")"
            printf "%s \tEmail: admin@nanocloud.com\n" "$(date "${DATE_FMT}")"
            printf "%s \tPassword: admin\n" "$(date "${DATE_FMT}")"
            echo "$(date "${DATE_FMT}") This URL will only be accessible from this host."
            echo ""
            echo "$(date "${DATE_FMT}") Use the following commands as root to start, stop or get status information"
            echo "$(date "${DATE_FMT}")     # installer.sh start"
            echo "$(date "${DATE_FMT}")     # installer.sh stop"
            echo "$(date "${DATE_FMT}")     # installer.sh status"
        fi
        ;;
    "uninstall")
        ./nanocloud/uninstall.sh
        rm -rf nanocloud
        ;;
esac

if [ -z "${COMMAND}" ]; then
    # If no argument is provided then install Nanocloud
    docker run -e HOST_UID=$SCRIPT_UID -v ${PWD}/nanocloud:/var/lib/nanocloud --rm nanocloud/community:${COMMUNITY_TAG}
    docker-compose -f ${PWD}/nanocloud/docker-compose.yml up -d
fi
