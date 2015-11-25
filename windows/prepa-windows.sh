#!/bin/bash

CPUS=4
RAM=2560
HD="output-windows-2012R2/windows-server-2012R2-amd64.qcow2"


QEMU=$(which qemu-system-x86_64)
${QEMU} \
-machine type=pc,accel=kvm \
-smp ${CPUS} \
-m ${RAM} \
-drive file=${HD},if=virtio,cache=writeback,discard=ignore \
-fda floppy/floppy_setup.img \
-fdb floppy/floppy_sysprep.img \
-vnc 0.0.0.0:87 \
-name packer-windows-2012R2 \
-netdev user,id=user.0,hostfwd=tcp::2222-:22 \
-device virtio-net,netdev=user.0 \
-redir tcp:3389::3389 \
-usb -device usb-tablet \
-monitor stdio
