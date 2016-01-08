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


DATE_FMT="+%Y/%m/%d %H:%M:%S"

NC_QEMU_PID=$(pgrep -fl nanocloud | awk '/qemu-system-x86/ { print $1; }')
echo "$(date "${DATE_FMT}") Stopping Nanocloud virtual machines"
for PID in $NC_QEMU_PID; do
    kill "${PID}"
    sleep 1
done

echo "$(date "${DATE_FMT}") Stopping host API"
/etc/init.d/iaasAPI stop > /dev/null 2>&1

# Check ip_forward
if [ "$(sysctl --value net.ipv4.ip_forward)" != "0" ]; then
    echo "$(date "${DATE_FMT}") Stoping IP Forward"
    sysctl --write net.ipv4.ip_forward=0 > /dev/null 2>&1
fi

echo "$(date "${DATE_FMT}") Nanocloud stopped"
echo "$(date "${DATE_FMT}") To start again Nanocloud, use :"
echo "$(date "${DATE_FMT}")     # /var/lib/nanocloud/scripts/start.sh"
