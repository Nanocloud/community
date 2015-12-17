#!/usr/bin/env bash
HERE=$(dirname $(readlink -f $0))

# symlink current Windows user dir into /home
ln -s "$(dirname $(cygpath -D))" /home/${USERNAME}

mkgroup -l > /etc/group
mkpasswd -l -p "$(cygpath -H)" > /etc/passwd

ssh-host-config -y --cygwin "ntsecbinmode mintty nodosfilewarning" --pwd "Nanocloud123+"

# Disable user / group permission checking
sed -i 's/.*StrictModes.*/StrictModes no/' /etc/sshd_config

# Disable reverse DNS lookups
sed -i 's/.*UseDNS.*/UseDNS no/' /etc/sshd_config
