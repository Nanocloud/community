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

${CURRENT_DIR}/installer/check_version.sh

# Power down all production container
docker-compose down

# Remove all docker images related to Nanocloud
docker images | awk '/^nanocloud\// { printf "docker rmi -f %s:%s\n", $1, $2; }' | sh

echo "$(date "${DATE_FMT}") Nanocloud uninstalled"
echo "$(date "${DATE_FMT}") To install Nanocloud again, use :"
echo "$(date "${DATE_FMT}")     # $PWD/install.sh"
