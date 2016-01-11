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

NANOCLOUD_DIR=${NANOCLOUD_DIR:-"${CURRENT_DIR}/installation_dir"}

if [ -z "$(which docker)" ]; then
  echo "$(date "${DATE_FMT}") Docker is missing, please install *docker*"
  exit 2
fi
if [ -z "$(which docker-compose)" ]; then
  echo "$(date "${DATE_FMT}") Docker-compose is missing, please install *docker-compose*"
  exit 2
fi

rm -f ${NANOCLOUD_DIR}/pid/windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64.pid

rm -rf ${CURRENT_DIR}/dockerfiles/build_output
docker-compose -f ${CURRENT_DIR}/dockerfiles/docker-compose.yml rm -f > /dev/null 2>&1
docker rmi -f guacamole-client \
       dockerfiles_guacamole-server \
       dockerfiles_nanocloud-api \
       dockerfiles_proxy \
       dockerfiles_ambassador \
       dockerfiles_rabbitmq \
       dockerfiles_postgres \
       dockerfiles_apps-module \
       dockerfiles_history-module \
       dockerfiles_iaas-module \
       dockerfiles_apiiaas-module \
       dockerfiles_ldap-module \
       dockerfiles_users-module > /dev/null 2>&1
docker rmi -f nanocloud/guacamole-client \
       nanocloud/guacamole-server \
       nanocloud/nanocloud-api \
       nanocloud/proxy \
       nanocloud/ambassador \
       nanocloud/rabbitmq \
       nanocloud/postgres \
       nanocloud/apps-module \
       nanocloud/history-module \
       nanocloud/iaas-module \
       nanocloud/apiiaas-module \
       nanocloud/ldap-module \
       nanocloud/users-module > /dev/null 2>&1

echo "$(date "${DATE_FMT}") Removing installed files"
