# Nanocloud community

Current version: **0.5.1**

Experience the seamless transformation of your application.

Nanocloud Community is an open source solution to provide a simple DevOps box
for techies and specialists. It allows every application to be pushed out in
any browser. A simple graphical interface allows users to experience any
applications with new uses.


## Quickinstall

Simply run the following command to install and run **Nanocloud**:

```
curl --progress-bar "http://releases.nanocloud.org:8080/releases/latest/installer.sh" | sh
```

> Read carefully outputs for this script, it will prodive you usefull
> information about the installation directory and URL where your installation
> will be accessible.

### Alternative to curl

If you don't want to or cannot use *curl*, you can launch the **installer** this
way:

```
wget "http://releases.nanocloud.org:8080/releases/latest/installer.sh" -q -O - | sh
```

## Prerequisites

For your host computer

* CPU must support hardware virtualization (Intel VT-x or AMD-V).
* Operating system must be a linux 64 bit. It is advised to use Debian 8 or
  later. Other Linux distributions may work.
* At least 4GB of RAM available.
* At least 7GB disk space (10GB recommended, depending on the software you want
  to deploy).

You also need to install the following packages on your distribution:

* *docker* 1.10+
* *docker-compose* 1.6.2
* *curl* to check status

## Uninstall

To uninstall Nanocloud, run the script nanocloud_uninstall.sh located in the root of the install directory

````
./nanocloud_uninstall.sh
````

Then you will need to manually remove your current directory

## How to build

To build Nanocloud Community involves building both Nanocloud and a Windows image ready to host your applications.

### Build Nanocloud

Nanocloud's components all run in their own dedicated Docker container.
Recipies to build them are defined in Dockerfiles and their launch is orchestrated with docker-compose

Simply run :

```
docker-compose -f modules/docker-compose-build.yml up -d
```

### Build Windows

At this point Nanocloud is running and right after login it is possible to download a built Windows image directly from the web interface.

If you are ready to wait 30 minutes to build a new image,you will need some packages :

* *packer*
* *qemu*
* *sshpass*
* *netcat*

Then run :

```
env PACKER_LOG=1 ./windows/build-windows.sh
```

Afterwards, You can run the following command to copy your Windows image to your running IaaS container.

```
docker cp windows/output-windows-2012R2-qemu/windows-server-2012R2-amd64.qcow2 iaas-module:/var/lib/nanocloud/images/windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64.qcow2
```

## Known bugs

**Nanocloud** is in active development phase. Some issues are known and
will be fixed in future releases.

If your issue isn't listed bellow, please report it
[here](https://github.com/Nanocloud/community/issues/new)

* While downloading **Windows**, information disappears from home page.

## Developer setup

Nanocloud developer environment is based on Docker containers. Development containers are based on production containers and add some features for developers such as :
- Code live reload on modification
- Debugger port mapped
- Software ran in debug mode

To run a dev environment, use the following command (after building Nanocloud as described above):

```
docker-compose -f modules/docker-compose-dev.yml up -d
```

You should be able to modify source on your local repository and see changes
directly applied.

You can access the log file with the following command (from the root directory):

```
docker-compose -f modules/docker-compose-dev.yml logs
```

## Roadmap

In future releases, we plan to add :

* Customize applications names and icons.
* Assign permission per application and per users.
* Dashboard to get information on the hosting platform.
* Buttons to log off a user from its windows session.
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
