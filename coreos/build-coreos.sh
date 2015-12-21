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


CHANNEL="beta"
COREOS_BASE_URL="http://${CHANNEL}.release.core-os.net/amd64-usr/current"
DATE_FMT="+%Y/%m/%d %H:%M:%S"


echo "$(date "${DATE_FMT}") Download CoreOS script…"
curl --progress-bar "${COREOS_BASE_URL}/coreos_production_qemu.sh" --output coreos_production_qemu.sh
echo "$(date "${DATE_FMT}") Download CoreOS image…"
curl --progress-bar "${COREOS_BASE_URL}/coreos_production_qemu_image.img.bz2" --output coreos_production_qemu_image.img.bz2

echo "$(date "${DATE_FMT}") Unpacking CoreOS…"
bzip2 -d coreos_production_qemu_image.img.bz2
chmod +x coreos_production_qemu.sh

echo "$(date "${DATE_FMT}") Generating SSH keys"
(
    echo -e "\n\n\n" | ssh-keygen -t rsa -N "" -f coreos.key
    chmod 400 coreos.key
) > /dev/null 2>&1

echo "$(date "${DATE_FMT}") Adding disk space to CoreOS…"
qemu-img resize coreos_production_qemu_image.img +5G

nohup ./coreos_production_qemu.sh -a coreos.key.pub -- -nographic > /dev/null &

echo "$(date "${DATE_FMT}") Testing connectivity…"
sleep 10
nc -v -z -w 10 localhost 2222 > /dev/null 2>&1
if [ "$?" != "0" ]; then
  echo "$(date "${DATE_FMT}") CoreOS failed to boot, exiting"
  exit 1
fi

echo "$(date "${DATE_FMT}") Provisioning…"
ssh \
    -o StrictHostKeyChecking=no \
    -i coreos.key \
    -l core \
    -p 2222 \
    localhost < "provision-coreos.sh"

echo "$(date "${DATE_FMT}") Compressing QCOW2 image…"
qemu-img convert -c -f qcow2 -O qcow2 coreos_production_qemu_image.img coreos.qcow2




