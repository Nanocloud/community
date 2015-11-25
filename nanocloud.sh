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


NANOCLOUD_DIR="/var/lib/nanocloud"
NANOCLOUD_BIN_URL="https://community.nanocloud.com/nanocloud"

COREOS_QCOW2_URL="http://stable.release.core-os.net/amd64-usr/current/coreos_production_qemu_image.img.bz2"
NANOCLOUD_QCOW2_URL="http://community.nanocloud.com/coreos.qcow2"
DATE_FMT="+%Y/%m/%d %H:%M:%S"

nano_exec() {
  # Arrange for the temporary file to be deleted when the script terminates
  trap 'rm -f "/tmp/exec.$$"' 0
  trap 'exit $?' 1 2 3 15

  # Create temporary file from the standard input
  cat >/tmp/exec.$$

  # Make the temporary file executable
  chmod +x /tmp/exec.$$

  # Execute the temporary file
  /tmp/exec.$$
}

# Check if current user is root
if [ "$(id -u)" != "0" ]; then
   echo "$(date "${DATE_FMT}") You must be root to run this script"
   exit 1
fi

echo "$(date "${DATE_FMT}") Activating *ip_forward*"
if [ "$(sysctl --value net.ipv4.ip_forward)" != "1" ]; then
  sysctl --write net.ipv4.ip_forward=1 > /dev/null 2>&1
fi

if [ "${1}" != "local" ]; then
  echo "$(date "${DATE_FMT}") Downloading Nanocloud binaries"
  wget --quiet --progress=bar:force ${NANOCLOUD_BIN_URL} -O - | nano_exec
  if [ "$?" != "0" ]; then
    echo "$(date "${DATE_FMT}") Installation failed, exiting…"
    exit 1
  fi
  echo "$(date "${DATE_FMT}") Downloading Coreos…"
  (
    cd ${NANOCLOUD_DIR}/images
    wget --quiet --progress=bar:force ${NANOCLOUD_QCOW2_URL} -O coreos.qcow2
    echo "$(date "${DATE_FMT}") Coreos download finished"
  )
else
  ./nanocloud unpack
  (
    cd ${NANOCLOUD_DIR}/images
    echo "$(date "${DATE_FMT}") Downloading CoreOS"
    wget --quiet --progress=bar:force ${COREOS_QCOW2_URL} -O - | bzcat > coreos.qcow2
  )
  ./nanocloud launch
fi

echo "$(date "${DATE_FMT}") Starting first VM…"
(
  cd ${NANOCLOUD_DIR}
  nohup scripts/launch-coreos.sh > start.log & 2>&1
)
chmod 400 /var/lib/nanocloud/coreos.key

echo "$(date "${DATE_FMT}") Testing connectivity…"
sleep 10
nc -v -z -w 10 localhost 2222 > /dev/null 2>&1
if [ "$?" != "0" ]; then
  echo "$(date "${DATE_FMT}") CoreOS failed to boot, exiting"
  exit 1
fi

echo "$(date "${DATE_FMT}") Setup complete"
echo "$(date "${DATE_FMT}") You can now manage your platform on : https://localhost:8443"
echo "$(date "${DATE_FMT}") Default admin credential:"
printf "%s \tEmail: admin@nanocloud.com\n" "$(date "${DATE_FMT}")"
printf "%s \tPassword: admin\n" "$(date "${DATE_FMT}")"
echo "$(date "${DATE_FMT}") This URL will only be accessible from this host."
echo ""
echo "$(date "${DATE_FMT}") Use the following commands as root to start, stop or get status information"
echo "$(date "${DATE_FMT}")     # /var/lib/nanocloud/scripts/start.sh"
echo "$(date "${DATE_FMT}")     # /var/lib/nanocloud/scripts/stop.sh"
echo "$(date "${DATE_FMT}")     # /var/lib/nanocloud/scripts/status.sh"
