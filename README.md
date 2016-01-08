# Nanocloud community

Current version: **0.1**

Experience the seamless transformation of your application.

Nanocloud Community is an open source solution to provide a simple DevOps box
for techies and specialists. It allows every application to be pushed out in
any browser. A simple graphical interface allows users to experience any
applications with new uses.


## Quickinstall

Simply run the following command as **root** to install and run **Nanocloud**:

```
curl --progress-bar "http://community.nanocloud.com/nanocloud.sh" | sh
```

> Note: You need to be *root* on the host machine to run **Nanocloud**. This
> won't be necessary any more in the next release.

At the end of the installation Nanocloud community will be installed in
*/var/lib/nanocloud*.

### Alternative to curl

If you don't want to or cannot use *curl*, you can launch the **one-liner** this way:

```
wget "http://community.nanocloud.com/nanocloud.sh" -q -O - | sh
```

## Prerequisites

For your host computer

* You must be able to login as root.
* CPU must support hardware virtualization (Intel VT-x or AMD-V).
* Operating system must be a linux 64 bit. It is advised to use Debian 8 or
  later. Other Linux distributions may work.
* At least 4GB of RAM available (6GB recommended).
* At least 6GB disk space (10GB recommended, depending on the software you want to
  deploy).

You also need to install the following packages on your distribution:

* *curl* or *wget*
* *netcat*

## How to build

You will need some tools to build Nanocloud

* gcc
* git
* go
* packer
* qemu-system-x86_64
* ssh_pass

Then to build your own installer, use this script:

```
./build.sh
```

It will build windows image, coreos image and then nanocloud tools. You can
launch the following command to install your build on your system:

```
./nanocloud.sh
```

> Note: for now, you have to be *root* to execute the lase command. This will change in next release

## Known bugs

**Nanocloud** is in active development phase. Some issues are known and
will be fixed in future releases.

If your issue isn't listed bellow, please report it
[here](https://github.com/Nanocloud/community/issues/new)

* This installation won't work with **parallels** or **virtualbox**.
* While downloading **Windows**, information disappears from home page.
* CoreOS qcow2 disk keeps growing until it's full and cause some 
[issue](http://stackoverflow.com/questions/31712266/how-to-clean-up-docker-overlay-directory).
* When CoreOS VM is stopped, users are erased.
* When an application is published, a connection for admin user is displayed but will not work.

## Roadmap

In future releases, we plan to add :

* Installation as a non-root user.
* Customize applications names and icons.
* Assign permission per application and per users.
* Dashboard to get information on the hosting platform.
* Buttons to log off a user from its windows session.
* A better authentication method.
* Display live metrics graphs (RAM/DISK/CPU usage).
* Change users information (not just password).
* Show last connection per user.

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
