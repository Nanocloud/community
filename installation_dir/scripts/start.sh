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


NANOCLOUD_DIR="/var/lib/nanocloud"
DATE_FMT="+%Y/%m/%d %H:%M:%S"


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

# Check ip_forward
echo "$(date "${DATE_FMT}") Activating *ip_forward*"
if [ "$(sysctl --value net.ipv4.ip_forward)" != "1" ]; then
  sysctl --write net.ipv4.ip_forward=1 > /dev/null 2>&1
fi

echo "$(date "${DATE_FMT}") Starting host API"
/etc/init.d/iaasAPI start > /dev/null 2>&1

(
  cd ${NANOCLOUD_DIR}
  nohup scripts/launch-coreos.sh > start.log &
)

echo "$(date "${DATE_FMT}") Nanocloud started"
printf "%s \tURL: https://localhost\n" "$(date "${DATE_FMT}")"
printf "%s \tEmail: admin@nanocloud.com\n" "$(date "${DATE_FMT}")"
printf "%s \tPassword: admin\n" "$(date "${DATE_FMT}")"
echo "$(date "${DATE_FMT}") This URL will only be accessible from this host."
echo ""
echo "$(date "${DATE_FMT}") Use the following commands as root to start, stop or get status information"
echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/start.sh"
echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/stop.sh"
echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/status.sh"
