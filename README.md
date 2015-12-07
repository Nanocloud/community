# Nanocloud community


## Quickinstall

Simply run the following command as **root** to install and run **Nanocloud**

```
curl --progress-bar "http://community.nanocloud.com/nanocloud.sh" | sh
```

> Note: You need to be *root* on the host machine to run **Nanocloud**. This
> will be improved in next release

At the end of the installation Nanocloud community will be installed in
*/var/lib/nanocloud*.

### Alternative to curl

If you don't want to or cannot use *curl*, you can launch the **one-liner** this way :

```
wget "http://community.nanocloud.com/nanocloud.sh" -q -O - | sh
```

## Prerequisites

For your host computer

* You must be able to login as root
* CPU must support harware virtualization (Intel VT-x or AMD-V).
* Operating system must be a linux 64 bits. It is advised to use Debian 8 or
  later. Other linux distributions may work.
* At least 4GB of RAM available (6GB recommended)
* At least 6GB disk space (10GB recommended, depending on software you want to
  deploy)

Then, you need to install the following package on your distribution

* *qemu-system-x86*
* *curl* or *wget* or *nc*


## How to build

If you want to build your own installer. Follow these step :

```
./build_nanocloud.sh
packer build windows/windows-2012-R2-standard-amd64.json
```

And, that's it.

## Licence

This file is part of Nanocloud community.

Nanocloud community is free software; you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

Nanocloud community is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
