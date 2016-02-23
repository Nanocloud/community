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

DATE_FMT="+%Y/%m/%d %H:%M:%S"

DOCKER_VERSION_REQUIRED_MAJOR=1
DOCKER_VERSION_REQUIRED_MINOR=10
DOCKER_VERSION_REQUIRED_FIX=0

DOCKER_COMPOSE_VERSION_REQUIRED_MAJOR=1
DOCKER_COMPOSE_VERSION_REQUIRED_MINOR=5
DOCKER_COMPOSE_VERSION_REQUIRED_FIX=2

check_docker_compose() {

    if [ -z "$(which docker-compose)" ]; then
        echo "$(date "${DATE_FMT}") Docker-compose is missing, please install *docker-compose*"
        return 0
    fi

    # Expected output for docker-compose version = MAJOR.MINOR.FIX
    DOCKER_COMPOSE_VERSION_MAJOR=$(docker-compose version --short | cut -d '.' -f 1)
    DOCKER_COMPOSE_VERSION_MINOR=$(docker-compose version --short | cut -d '.' -f 2)
    DOCKER_COMPOSE_VERSION_FIX=$(docker-compose version --short | cut -d '.' -f 3)

    # For now Nanocloud is only compatible with docker-compose 1.5.2
    if [ $DOCKER_COMPOSE_VERSION_MAJOR -ne $DOCKER_COMPOSE_VERSION_REQUIRED_MAJOR ]; then
        return 0
    fi
    if [ $DOCKER_COMPOSE_VERSION_MINOR -ne $DOCKER_COMPOSE_VERSION_REQUIRED_MINOR ]; then
        return 0
    fi
    if [ $DOCKER_COMPOSE_VERSION_FIX -ne $DOCKER_COMPOSE_VERSION_REQUIRED_FIX ]; then
        return 0
    fi

    return 1
}

check_docker() {

    if [ -z "$(which docker)" ]; then
        echo "$(date "${DATE_FMT}") Docker is missing, please install *docker*"
        return 0
    fi

    # Expected output for docker version = MAJOR.MINOR.FIX
    DOCKER_VERSION_MAJOR=$(docker version --format '{{.Server.Version}}' | cut -d '.' -f 1)
    DOCKER_VERSION_MINOR=$(docker version --format '{{.Server.Version}}' | cut -d '.' -f 2)
    DOCKER_VERSION_FIX=$(docker version --format '{{.Server.Version}}' | cut -d '.' -f 3)

    # Test major version
    if [ $DOCKER_VERSION_MAJOR -eq $DOCKER_VERSION_REQUIRED_MAJOR ]; then

        # Test minor version
        if [ $DOCKER_VERSION_MINOR -ge $DOCKER_VERSION_REQUIRED_MINOR ]; then

            # Good if version is above required version
            if [ $DOCKER_VERSION_MINOR -ne $DOCKER_VERSION_REQUIRED_MINOR ]; then
                return 1
            fi

            # Test fix version
            if [ $DOCKER_VERSION_FIX -ge $DOCKER_VERSION_REQUIRED_FIX ]; then
                return 1
            fi
        fi
    fi

    return 0
}

check_dependencies() {

    if check_docker_compose -eq 1 ; then
        echo "$(date "${DATE_FMT}") Installed docker-compose is incompatible. Expected " $DOCKER_COMPOSE_VERSION_REQUIRED_MAJOR.$DOCKER_COMPOSE_VERSION_REQUIRED_MINOR.$DOCKER_COMPOSE_VERSION_REQUIRED_FIX " but found " $(docker-compose version --short)
        exit 2
    fi

    if check_docker -eq 1 ; then
        echo "$(date "${DATE_FMT}") Installed docker is incompatible. Expected " $DOCKER_VERSION_REQUIRED_MAJOR.$DOCKER_VERSION_REQUIRED_MINOR.$DOCKER_VERSION_REQUIRED_FIX " but found " $(docker version --format '{{.Server.Version}}')
        exit 2
    fi
}

check_dependencies
