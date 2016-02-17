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
DATE_FMT="+%Y/%m/%d %H:%M:%S"

WINDOWS_QCOW2_FILENAME="${CURRENT_DIR}/output-windows-2012R2-qemu/windows-server-2012R2-amd64.qcow2"
WINDOWS_PASSWORD="Nanocloud123+"
VM_HOSTNAME="windows-2012R2"
VM_NCPUS="$(grep -c ^processor /proc/cpuinfo)"
SSH_PORT=1119
QEMU=$(which qemu-system-x86_64 || true)

if [ -z "$(which packer || true)" ]; then
  echo "$(date "${DATE_FMT}") Packer is missing, please install *packer*"
  exit 2
fi
if [ -z "${QEMU}" ]; then
  echo "$(date "${DATE_FMT}") Qemu is missing, please install *qemu*"
  exit 2
fi
if [ -z "$(which sshpass || true)" ]; then
  echo "$(date "${DATE_FMT}") sshpass is missing, please install *sshpass*"
  exit 2
fi
if [ -z "$(which netcat || true)" ]; then
  echo "$(date "${DATE_FMT}") netcat is missing, please install *netcat*"
  exit 2
fi

if [ ! -f "${WINDOWS_QCOW2_FILENAME}" ]; then
    (
        cd "${CURRENT_DIR}"
        packer build -only=windows-2012R2-qemu windows_2012_r2.json
    )
fi

nohup "${QEMU}" \
    -nodefaults \
    -name "${VM_NAME}" \
    -m 4096 \
    -cpu host \
    -smp "${VM_NCPUS}" \
    -machine accel=kvm \
    -drive if=virtio,file="${WINDOWS_QCOW2_FILENAME}" \
    -vnc :2 \
    -pidfile "${CURRENT_DIR}/${VM_HOSTNAME}.pid" \
    -net nic,vlan=0,model=virtio \
    -net user,vlan=0,hostfwd=tcp::"${SSH_PORT}"-:22,hostname="${VM_HOSTNAME}" \
    -vga qxl \
    -global qxl-vga.vram_size=33554432 \
    "${@}" &

echo "$(date "${DATE_FMT}") Waiting 50 seconds for windows to boot"
sleep 50

echo "$(date "${DATE_FMT}") Testing connectivity…"
nc -v -z -w 10 localhost ${SSH_PORT} > /dev/null 2>&1
if [ "$?" != "0" ]; then
  echo "$(date "${DATE_FMT}") Windows failed to boot, exiting"
  exit 1
fi

echo "$(date "${DATE_FMT}") Installing new Remote Desktop session deployment…"
sshpass -p "${WINDOWS_PASSWORD}" ssh -p ${SSH_PORT} -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null Administrator@localhost << EOF
powershell.exe "import-module remotedesktop ; New-RDSessionDeployment -ConnectionBroker adapps.intra.localdomain.com -WebAccessServer adapps.intra.localdomain.com -SessionHost adapps.intra.localdomain.com; New-RDSessionCollection -CollectionName collection -SessionHost adapps.intra.localdomain.com -CollectionDescription 'Nanocloud collection' -ConnectionBroker adapps.intra.localdomain.com; New-RDRemoteApp -CollectionName collection -DisplayName hapticPowershell -FilePath 'C:\Windows\system32\WindowsPowerShell\v1.0\powershell.exe' -Alias hapticPowershell -CommandLineSetting Require -RequiredCommandLine '-ExecutionPolicy Bypass c:\publishApplication.ps1'"
powershell.exe "sc.exe config RDMS start= auto"
EOF

echo "$(date "${DATE_FMT}") Pushing hapticPowershell script to Windows…"
sshpass -p "${WINDOWS_PASSWORD}" scp -P ${SSH_PORT} -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null "${CURRENT_DIR}/floppy/windows-2012-standard-amd64/publishApplication.ps1" Administrator@localhost:"C:\publishApplication.ps1"

sshpass -p "${WINDOWS_PASSWORD}" ssh -p ${SSH_PORT} -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null Administrator@localhost << EOF
powershell.exe "shutdown.exe /s /f /d p:4:1 /c 'Provisioning Shutdown'"
EOF

echo "$(date "${DATE_FMT}") Waiting for windows to shutdown…"
sleep 20

echo "$(date "${DATE_FMT}") Compressing QCOW2 image…"
qemu-img convert -c -f qcow2 -O qcow2 "${WINDOWS_QCOW2_FILENAME}" "${WINDOWS_QCOW2_FILENAME}.mini"
mv "${WINDOWS_QCOW2_FILENAME}.mini" "${WINDOWS_QCOW2_FILENAME}"
