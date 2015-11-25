#!/bin/bash -x

VM_NAME="winad-milli-free_use-10.20.12.20-windows-server-std-2012-x86_64"
DIR_NAME=$(dirname "${0}")
SCRIPT_DIR=$(readlink -e "${DIR_NAME}")
NANOCLOUD_DIR="/var/lib/nanocloud/power"
SYSTEM_VHD="${SCRIPT_DIR}/output-windows-2012R2/windows-server-2012R2-amd64.qcow2"

ETH0=$(cat /proc/net/arp | grep -v tap | grep -vi Device | awk '{print $6}' | uniq | sort -n | head -n 1)
PUBLIC_IP=$(ip -4 addr show dev ${ETH0} | grep inet | grep -v ${ETH0}: | awk '{print $2}' | sed -s 's/\/.*//' | head -n 1)
VM_IP=10.20.12.20
VM_INTERFACE=tap-020-012-020

if [ ! -f "${SYSTEM_VHD}" ]; then
    echo "You must build packer template before running this script"
    exit 1
fi

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
$QEMU \
    -nodefaults \
    -name ${VM_NAME} \
    -pidfile ${NANOCLOUD_DIR}/pid/${VM_NAME}.pid  \
    -drive file=${SYSTEM_VHD},if=none,media=disk,cache=writeback,aio=native,id=vhd_system \
    -device driver=virtio-blk-pci,drive=vhd_system \
    -drive file='',if=ide,media=cdrom \
    -fda floppy/floppy_setup.img \
    -boot order=cd \
    -machine type=pc,accel=kvm \
    -cpu host,kvm=off \
    -smp 4 \
    -m 2580 \
    -usb \
    -device usb-tablet \
    -rtc base=localtime,clock=host \
    -chardev socket,id=monitor,path=${NANOCLOUD_DIR}/socket/${VM_NAME}.socket,server,nowait \
    -mon chardev=monitor,mode=readline \
    -netdev type=tap,id=tap_id_19539f4,ifname=${VM_INTERFACE},script=no,downscript=no,vhost=on \
    -device virtio-net-pci,netdev=tap_id_19539f4,mac=52:54:00:53:af:3a \
    -soundhw hda \
    -vnc 0.0.0.0:87 \
    -vga qxl \
    -nographic \
    -global qxl-vga.vram_size=33554432 \
    -spice port=8097,addr=${PUBLIC_IP},password='firstpass',streaming-video=off \
    -device virtio-serial \
    -chardev spicevmc,id=vdagent,name=vdagent \
    -device virtserialport,chardev=vdagent,name=com.redhat.spice.0 \
    -chardev spicevmc,name=usbredir,id=usbredirchardev1 \
    -device usb-redir,chardev=usbredirchardev1,id=usbredirdev1 \
    -chardev spicevmc,name=usbredir,id=usbredirchardev2 \
    -device usb-redir,chardev=usbredirchardev2,id=usbredirdev2 \
    -chardev spicevmc,name=usbredir,id=usbredirchardev3 \
    -device usb-redir,chardev=usbredirchardev3,id=usbredirdev3 \
    -watchdog i6300esb \
    -watchdog-action none \
    -uuid 19539f4c-d88d-40e9-9cba-d240c2f1ad8d \
    -balloon virtio \
    -monitor stdio

#    -fda file="/root/packer/packer-qemu-templates/windows/Autounattend.fd" \

#    -device vfio-pci,host=00:02.0,x-vga=on \
#    -device vfio-pci,host=00:03.0 \

#/var/lib/nanocloud/qemu/bin/qemu-system-x86_64 -drive file=output-windows-2012R2/packer-windows-2012R2,if=virtio,cache=writeback,discard=ignore -vnc 0.0.0.0:87 -name packer-windows-2012R2 -smp 4 -device virtio-net-pci,netdev=tap_id_19539f4,mac=52:54:00:53:af:3a -boot once=d -machine type=pc,accel=kvm -m 2560 -netdev type=tap,id=tap_id_19539f4,ifname=tap-020-012-020,script=no,downscript=no,vhost=on -cdrom /root/packer/packer-qemu-templates/windows/packer_cache/5528757477328a4f4cef4cfd33a8b48d3765375bcbbf18340434a6918cd34740.ISO



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

TFINAL=$(($($DATE +%s)-$DTEPOCH))
echo "Runtime=$(($($DATE +%s)-$DTEPOCH)) s"

/bin/rm ${NANOCLOUD_DIR}/pid/winad-milli-free_use-10.20.12.20-windows-server-std-2012-x86_64.pid
/bin/rm ${NANOCLOUD_DIR}/socket/winad-milli-free_use-10.20.12.20-windows-server-std-2012-x86_64.socket
