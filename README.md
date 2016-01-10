# Nanocloud community

Current version: **0.2**

Experience the seamless transformation of your application.

Nanocloud Community is an open source solution to provide a simple DevOps box
for techies and specialists. It allows every application to be pushed out in
any browser. A simple graphical interface allows users to experience any
applications with new uses.


## Quickinstall

Simply run the following command to install and run **Nanocloud**:

```
curl --progress-bar "http://community.nanocloud.com/installer.sh" | sh
```

> Read carefully outputs for this script, it will prodive you usefull
> information about the installation directory and URL where your installation
> will be accessible.

### Alternative to curl

If you don't want to or cannot use *curl*, you can launch the **installer** this
way:

```
wget "http://community.nanocloud.com/installer.sh" -q -O - | sh
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

* *docker*
* *curl* or *wget*
* *netcat*

## How to build

Then to build your own installer, use this script:

```
./build.sh
```

It will build windows image and then nanocloud containers.
Afterwards, You can launch the following command to install your build on your
system:

```
./nanocloud.sh
```

## Known bugs

**Nanocloud** is in active development phase. Some issues are known and
will be fixed in future releases.

If your issue isn't listed bellow, please report it
[here](https://github.com/Nanocloud/community/issues/new)

* While downloading **Windows**, information disappears from home page.

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
