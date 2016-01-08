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


# Check if current user is root
NANOCLOUD_DIR="/var/lib/nanocloud"
DATE_FMT="+%Y/%m/%d %H:%M:%S"
NC_QEMU_PID=$(pgrep -fl nanocloud | awk '/qemu-system-x86/ { print $1; }')

# Check ip_forward
if [ "$(sysctl --value net.ipv4.ip_forward)" != "1" ]; then
  echo "$(date "${DATE_FMT}") IP Forward is missing, please use the following command to fix it"
  echo "$(date "${DATE_FMT}")    # sysctl --write net.ipv4.ip_forward=1"
CURL_CMD=$(which curl)
WGET_CMD=$(which wget)
if [ -n "${CURL_CMD}" ]; then
    NANOCLOUD_STATUS=$(curl --output /dev/null --insecure --silent --write-out '%{http_code}\n' "https://localhost")
elif [ -n "${WGET_CMD}" ]; then
    NANOCLOUD_STATUS=$(LANG=C wget --no-check-certificate "https://localhost" -O /dev/null 2>&1 | awk '/^HTTP/ { print $6 ;}')
fi

if [ -z "$NC_QEMU_PID" ]; then
  echo "$(date "${DATE_FMT}") Nanocloud is *NOT* running"
else
  echo "$(date "${DATE_FMT}") Nanocloud is running"
  printf "%s \tURL: https://localhost\n" "$(date "${DATE_FMT}")"
  printf "%s \tEmail: admin@nanocloud.com\n" "$(date "${DATE_FMT}")"
  printf "%s \tPassword: admin\n" "$(date "${DATE_FMT}")"
  echo "$(date "${DATE_FMT}") This URL will only be accessible from this host."
  echo ""
  echo "$(date "${DATE_FMT}") Use the following commands as root to start, stop or get status information"
  echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/start.sh"
  echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/stop.sh"
  echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/status.sh"
fi
