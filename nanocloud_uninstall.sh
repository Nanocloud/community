#!/bin/bash
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

NANOCLOUD_DIR=${NANOCLOUD_DIR:-"/var/lib/nanocloud"}
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
echo "# Erasing previous install"
echo "$(date "${DATE_FMT}") Killing Nanocloud qemu processes"
NC_QEMU_PID=$(pgrep -fl nanocloud | grep qemu-system-x86 | awk '{ print $1; }')
for PID in $NC_QEMU_PID; do
    kill "${PID}"
    sleep 1
done

echo "$(date "${DATE_FMT}") Stoping API"
/etc/init.d/iaasAPI stop > /dev/null 2>&1

echo "$(date "${DATE_FMT}") Removing installed files"
[ -d "${NANOCLOUD_DIR}" ] && rm -rf "${NANOCLOUD_DIR}"
[ -h /etc/init.d/iaasAPI ] && rm -rf /etc/init.d/iaasAPI
