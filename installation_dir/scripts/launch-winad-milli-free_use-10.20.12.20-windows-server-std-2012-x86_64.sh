#!/bin/bash -x
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

VM_NAME="windows-server-2012R2-amd64"
VM_HOSTNAME="windows-server-2012R2-amd64"

# Port map
SSH_PORT=1119
RDP_PORT=3389
LDAPS_PORT=6360

QEMU=$(which qemu-system-x86_64)
SYSTEM_VHD="${NANOCLOUD_DIR}/images/${VM_NAME}.qcow2"
VM_NCPUS="$(grep -c ^processor /proc/cpuinfo)"

$QEMU \
    -nodefaults \
    -name "${VM_NAME}" \
    -m 2560 \
    -cpu host \
    -smp "${VM_NCPUS}" \
    -machine accel=kvm \
    -drive if=virtio,file="${SYSTEM_VHD}" \
    -vnc :2 \
    -pidfile ${NANOCLOUD_DIR}/pid/${VM_NAME}.pid  \
    -chardev socket,id=monitor,path=${NANOCLOUD_DIR}/sockets/${VM_NAME}.socket,server,nowait \
    -net nic,vlan=0,model=virtio \
    -net user,vlan=0,hostfwd=tcp::"${SSH_PORT}"-:22,hostfwd=tcp::"${RDP_PORT}"-:3389,hostfwd=tcp::"${LDAPS_PORT}"-:636,hostname="${VM_HOSTNAME}" \
    -vga qxl \
    -global qxl-vga.vram_size=33554432 \
    "${@}"

/bin/rm ${NANOCLOUD_DIR}/pid/${VM_NAME}.pid
/bin/rm ${NANOCLOUD_DIR}/sockets/${VM_NAME}.socket
