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

SCRIPT_UID=$(id -u)

COMMUNITY_TAG="0.2.1"
COMMAND=${1}

if [ "${COMMAND}" = "indiana" ]; then
    COMMUNITY_TAG="indiana"
fi

if [ -z "$(which docker || true)" ]; then
  echo "$(date "${DATE_FMT}") Docker is missing, please install *docker*"
  exit 2
fi
if [ -z "$(which docker-compose || true)" ]; then
  echo "$(date "${DATE_FMT}") Docker-compose is missing, please install *docker-compose*"
  exit 2
fi

docker run -e HOST_UID=$SCRIPT_UID -v ${PWD}/nanocloud:/var/lib/nanocloud nanocloud/community:${COMMUNITY_TAG}
${PWD}/nanocloud/installation_dir/scripts/start.sh ${COMMUNITY_TAG}
