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


SCRIPT_FULL_PATH=$(readlink -e "${0}")
CURRENT_DIR=$(dirname "${SCRIPT_FULL_PATH}")

NANOCLOUD_DIR=${NANOCLOUD_DIR:-"/var/lib/nanocloud"}
NANOCLOUD_BIN_URL=${NANOCLOUD_DIR:-"https://community.nanocloud.com/nanocloud"}

NANOCLOUD_REPOS_URL="http://releases.nanocloud.org:8080/indiana"
DATE_FMT="+%Y/%m/%d %H:%M:%S"

download() {
  CURL_CMD=$(which curl)
  WGET_CMD=$(which wget)

  URL=${1}
  if [ -n "${CURL_CMD}" ]; then
    curl --progress-bar "${URL}"
  elif [ -n "${WGET_CMD}" ]; then
    wget --quiet "${URL}" -O -
  else
    echo "You need *curl* or *wget* to run this script, exiting"
    exit 2
  fi
}

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

if [ -z "$(which qemu-system-x86_64)" ]; then
  echo "$(date "${DATE_FMT}") Qemu is missing, please install *qemu-system-x86_64*"
  exit 2
fi
if [ -z "$(which curl)" -o -z "$(which wget)" ]; then
  echo "$(date "${DATE_FMT}") No download method found, please install *curl* or *wget*"
  exit 2
fi
if [ -z "$(which nc)" -o -z "$(which netcat)" ]; then
  echo "$(date "${DATE_FMT}") Netcat not found, please install *nc* or *netcat* command"
  exit 2
fi

echo "$(date "${DATE_FMT}") Activating *ip_forward*"
if [ "$(sysctl --value net.ipv4.ip_forward)" != "1" ]; then
  sysctl --write net.ipv4.ip_forward=1 > /dev/null 2>&1
fi

NANOCLOUD_BINARIES="${CURRENT_DIR}/nanocloud"
if [ -f "${NANOCLOUD_BINARIES}" ]; then
    echo "$(date "${DATE_FMT}") Local Nanocloud binaries available"
    "${CURRENT_DIR}/nanocloud"
else
    echo "$(date "${DATE_FMT}") Downloading Nanocloud binaries"
    download "${NANOCLOUD_BIN_URL}" | nano_exec
    if [ "$?" != "0" ]; then
      echo "$(date "${DATE_FMT}") Installation failed, exiting…"
      exit 1
    fi
fi

COREOS_QCOW2_FILENAME="${CURRENT_DIR}/coreos/coreos.qcow2"
if [ -f "${COREOS_QCOW2_FILENAME}" ]; then
    echo "$(date "${DATE_FMT}") Local CoreOS disk available"
    cp "${COREOS_QCOW2_FILENAME}" "${NANOCLOUD_DIR}/images/coreos-custom-free_use-127.0.0.1-linux-coreos-x86_64.qcow2"
    cp "${CURRENT_DIR}/coreos/coreos.key" "${NANOCLOUD_DIR}/coreos.key"
    cp "${CURRENT_DIR}/coreos/coreos.key.pub" "${NANOCLOUD_DIR}/coreos.key.pub"
else
    echo "$(date "${DATE_FMT}") Downloading Coreos…"
    (
      cd "${NANOCLOUD_DIR}/images"
      download "${NANOCLOUD_REPOS_URL}/coreos.qcow2" > coreos-custom-free_use-127.0.0.1-linux-coreos-x86_64.qcow2
      echo "$(date "${DATE_FMT}") Coreos download finished"
    )
    (
      cd "${NANOCLOUD_DIR}"
      download "${NANOCLOUD_REPOS_URL}/coreos.key" > coreos.key
      download "${NANOCLOUD_REPOS_URL}/coreos.key.pub" > coreos.key.pub
      echo "$(date "${DATE_FMT}") Coreos download finished"
    )
fi

echo "$(date "${DATE_FMT}") Starting first VM…"
(
  cd "${NANOCLOUD_DIR}"
  nohup scripts/launch-coreos-custom-free_use-127.0.0.1-linux-coreos-x86_64.sh > start.log & 2>&1
)
chmod 400 "${NANOCLOUD_DIR}/coreos.key"

echo "$(date "${DATE_FMT}") Testing connectivity…"
sleep 10
nc -v -z -w 10 localhost 2222 > /dev/null 2>&1
if [ "$?" != "0" ]; then
  echo "$(date "${DATE_FMT}") CoreOS failed to boot, exiting"
  exit 1
fi

WINDOWS_QCOW2_FILENAME="${CURRENT_DIR}/windows/output-windows-2012R2-qemu/windows-server-2012R2-amd64.qcow2"
if [ -f "${WINDOWS_QCOW2_FILENAME}" ]; then
  echo "$(date "${DATE_FMT}") Local Windows image found, copying"
  cp "${WINDOWS_QCOW2_FILENAME}" "${NANOCLOUD_DIR}/images/windows-server-2012R2-amd64.qcow2"
fi

echo "$(date "${DATE_FMT}") Setup complete"
echo "$(date "${DATE_FMT}") You can now manage your platform on : https://localhost:8443"
echo "$(date "${DATE_FMT}") Default admin credential:"
printf "%s \tEmail: admin@nanocloud.com\n" "$(date "${DATE_FMT}")"
printf "%s \tPassword: admin\n" "$(date "${DATE_FMT}")"
echo "$(date "${DATE_FMT}") This URL will only be accessible from this host."
echo ""
echo "$(date "${DATE_FMT}") Use the following commands as root to start, stop or get status information"
echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/start.sh"
echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/stop.sh"
echo "$(date "${DATE_FMT}")     # ${NANOCLOUD_DIR}/scripts/status.sh"
