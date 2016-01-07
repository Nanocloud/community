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

VM_NAME="winad-milli-free_use-10.20.12.20-windows-server-std-2012-x86_64"
NANOCLOUD_DIR="/var/lib/nanocloud"

ETH0=$(cat /proc/net/arp | grep -v tap | grep -vi Device | awk '{print $6}' | uniq | sort -n | head -n 1)
PUBLIC_IP=$(ip -4 addr show dev ${ETH0} | grep inet | grep -v ${ETH0}: | awk '{print $2}' | sed -s 's/\/.*//' | head -n 1)
VM_IP=10.20.12.20
VM_INTERFACE=tap-020-012-020

ulimit -c $(ulimit -Hc)
ulimit -m $(ulimit -Hm)
ulimit -s $(ulimit -Hs)

DATE='/bin/date -u'
HEURE=$($DATE +%T)
JOUR=$($DATE +%F)
T0=$JOUR" "$HEURE

echo "$DAY $HOUR GMT: Nanocloud Virtual Server: ${VM_NAME} - starts"
DTEPOCH=$($DATE --date="$T0" +%s)

#Tap interface creation
echo "Configure ${VM_INTERFACE} interface ..."
tunctl -b -u root -t ${VM_INTERFACE}
ip link set dev ${VM_INTERFACE} addr 00:c6:36:fc:07:1b
echo "Activating link ${VM_INTERFACE} ..."
ip link set dev ${VM_INTERFACE} up

echo "Gateway IP: 10.20.12.19 on ${VM_INTERFACE} ..."
ip addr add 10.20.12.19/30 brd 10.20.12.21 dev ${VM_INTERFACE}
echo "Done."

echo "Activating routing on VM IP: ${VM_IP}"
ip route add ${VM_IP}/32 via 10.20.12.19 dev ${VM_INTERFACE}

# Iptables rules for each interface
iptables -A FORWARD -i ${VM_INTERFACE} -j ACCEPT
iptables -A FORWARD -o ${VM_INTERFACE} -j ACCEPT


# Masquerading towards outside
iptables -t nat -A POSTROUTING -s ${VM_IP} -j SNAT --to-source ${PUBLIC_IP}

# Iptables rules for VNC & SPICE 
iptables -t filter -A INPUT -p tcp -d ${PUBLIC_IP}  --dport 6997  -j ACCEPT
iptables -t filter -A INPUT -p tcp -d ${PUBLIC_IP}  --dport 8097  -j ACCEPT

# SSH
iptables -t filter -A INPUT -p tcp -d ${PUBLIC_IP} --dport 1119 -j ACCEPT
iptables -t nat -A PREROUTING -p tcp -d ${PUBLIC_IP} --dport 1119 -j DNAT --to ${VM_IP}:22
iptables -t filter -A FORWARD -p tcp -d ${VM_IP} --dport 22

# HTTP
iptables -t filter -A INPUT -p tcp -d ${PUBLIC_IP} --dport 1177 -j ACCEPT
iptables -t nat -A PREROUTING -p tcp -d ${PUBLIC_IP} --dport 1177 -j DNAT --to ${VM_IP}:80
iptables -t filter -A FORWARD -p tcp -d ${VM_IP} --dport 80

# HTTPS
iptables -t filter -A INPUT -p tcp -d ${PUBLIC_IP} --dport 1540 -j ACCEPT
iptables -t nat -A PREROUTING -p tcp -d ${PUBLIC_IP} --dport 1540 -j DNAT --to ${VM_IP}:443
iptables -t filter -A FORWARD -p tcp -d ${VM_IP} --dport 443

#RDP
iptables -t filter -A INPUT -p tcp -d ${PUBLIC_IP} --dport 3389 -j ACCEPT
iptables -t nat -A PREROUTING -p tcp -d ${PUBLIC_IP} --dport 3389 -j DNAT --to ${VM_IP}:3389
iptables -t filter -A FORWARD -p tcp -d ${VM_IP} --dport 3389

# LDAPS
iptables -t filter -A INPUT -p tcp -d ${PUBLIC_IP} --dport 636 -j ACCEPT
iptables -t nat -A PREROUTING -p tcp -d ${PUBLIC_IP} --dport 636 -j DNAT --to ${VM_IP}:636
iptables -t filter -A FORWARD -p tcp -d ${VM_IP} --dport 636

QEMU=$(which qemu-system-x86_64)

SYSTEM_VHD="${NANOCLOUD_DIR}/images/${VM_NAME}.qcow2"
VM_NCPUS="$(grep -c ^processor /proc/cpuinfo)"

$QEMU \
    -nodefaults \
    -name ${VM_NAME} \
    -enable-kvm \
    -cpu host \
    -smp "${VM_NCPUS}" \
    -m 2560 \
    -pidfile ${NANOCLOUD_DIR}/pid/${VM_NAME}.pid  \
    -drive file=${SYSTEM_VHD},if=none,media=disk,cache=writeback,aio=native,id=vhd_system \
    -device driver=virtio-blk-pci,drive=vhd_system \
    -rtc base=localtime,clock=host \
    -chardev socket,id=monitor,path=${NANOCLOUD_DIR}/sockets/${VM_NAME}.socket,server,nowait \
    -mon chardev=monitor,mode=readline \
    -netdev type=tap,id=tap_id_19539f4,ifname=${VM_INTERFACE},script=no,downscript=no,vhost=on \
    -device virtio-net-pci,netdev=tap_id_19539f4,mac=52:54:00:53:af:3a \
    -vnc :2 \
    -vga qxl \
    -global qxl-vga.vram_size=33554432 \
    -nographic


# Stops masquerading towards outside
iptables -t nat -D POSTROUTING -s ${VM_IP} -j SNAT --to-source ${PUBLIC_IP}

# VNC & SPICE
iptables -t filter -D INPUT -p tcp -d ${PUBLIC_IP}  --dport 6997  -j ACCEPT
iptables -t filter -D INPUT -p tcp -d ${PUBLIC_IP}  --dport 8097  -j ACCEPT

# SSH
iptables -t filter -D INPUT -p tcp -d ${PUBLIC_IP} --dport 1119 -j ACCEPT
iptables -t nat -D PREROUTING -p tcp -d ${PUBLIC_IP} --dport 1119 -j DNAT --to ${VM_IP}:22
iptables -t filter -D FORWARD -p tcp -d ${VM_IP} --dport 22

# HTTP
iptables -t filter -D INPUT -p tcp -d ${PUBLIC_IP} --dport 1177 -j ACCEPT
iptables -t nat -D PREROUTING -p tcp -d ${PUBLIC_IP} --dport 1177 -j DNAT --to ${VM_IP}:80
iptables -t filter -D FORWARD -p tcp -d ${VM_IP} --dport 80

# HTTPS
iptables -t filter -D INPUT -p tcp -d ${PUBLIC_IP} --dport 1540 -j ACCEPT
iptables -t nat -D PREROUTING -p tcp -d ${PUBLIC_IP} --dport 1540 -j DNAT --to ${VM_IP}:443
iptables -t filter -D FORWARD -p tcp -d ${VM_IP} --dport 443

# RDP
iptables -t filter -D INPUT -p tcp -d ${PUBLIC_IP} --dport 3389 -j ACCEPT
iptables -t nat -D PREROUTING -p tcp -d ${PUBLIC_IP} --dport 3389 -j DNAT --to ${VM_IP}:3389
iptables -t filter -D FORWARD -p tcp -d ${VM_IP} --dport 3389

# LDAPS
iptables -t filter -D INPUT -p tcp -d ${PUBLIC_IP} --dport 636 -j ACCEPT
iptables -t nat -D PREROUTING -p tcp -d ${PUBLIC_IP} --dport 636 -j DNAT --to ${VM_IP}:636
iptables -t filter -D FORWARD -p tcp -d ${VM_IP} --dport 636


# Stops interface
ip link set dev ${VM_INTERFACE} down

# Iptables rules for each interface
iptables -D FORWARD -i ${VM_INTERFACE} -j ACCEPT
iptables -D FORWARD -o ${VM_INTERFACE} -j ACCEPT

# Destroys interface
tunctl -d ${VM_INTERFACE}


DATE='/bin/date -u'
HEURE=$($DATE +%T)
JOUR=$($DATE +%F)

echo "$DAY $HOUR GMT: Nanocloud Virtual Server - ${VM_NAME} - stopped"

echo "Runtime=$(($($DATE +%s)-$DTEPOCH)) s"

/bin/rm ${NANOCLOUD_DIR}/pid/${VM_NAME}.pid
/bin/rm ${NANOCLOUD_DIR}/sockets/${VM_NAME}.socket
