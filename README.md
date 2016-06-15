# Nanocloud community [![Build Status](https://travis-ci.org/Nanocloud/community.svg?branch=master)](https://travis-ci.org/Nanocloud/community)

Current version: **0.7.0**

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

To uninstall Nanocloud, run the script uninstall.sh located in the root of the install directory

````
./uninstall.sh
````

Then you will need to manually remove your current directory

## How to build

To build Nanocloud Community involves building both Nanocloud and a Windows image ready to host your applications.

### Build Nanocloud

Nanocloud's components all run in their own dedicated Docker container.
Recipies to build them are defined in Dockerfiles and their launch is orchestrated with docker-compose

Simply run :

```
docker-compose -f modules/docker-compose-build.yml build
docker-compose -f modules/docker-compose-build.yml up -d
```

### Build Windows

At this point Nanocloud is running and right after login it is possible to download a built Windows image directly from the web interface.

To build an image from scratch is not yet documented since the procedure changed. https://github.com/Nanocloud/community/issues/502

However, you can download a fresh Windows image from Nanocloud's release server.

```
wget http://releases.nanocloud.org:8080/releases/latest/windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64.qcow2
docker cp windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64.qcow2 iaas-module:/var/lib/nanocloud/images/windows-custom-server-127.0.0.1-windows-server-std-2012R2-amd64.qcow2
```

## Developer setup

Nanocloud developer environment is based on Docker containers. Development containers are based on production containers and add some features for developers such as :
- Code live reload on modification
- Debugger port mapped
- Software ran in debug mode

To run a dev environment, use the following command (after building Nanocloud as described above):

```
docker-compose -f modules/docker-compose-build.yml build
docker-compose -f modules/docker-compose-dev.yml build
docker-compose -f modules/docker-compose-dev.yml up -d
```

You should be able to modify source on your local repository and see changes
directly applied.

You can access the log file with the following command (from the root directory):

```
docker-compose -f modules/docker-compose-dev.yml logs
```

## Configuration

You can configure *nanocloud* with the following environement variables in *modules/docker-compose.yml* for nanocloud-backend:

* ADMIN_FIRSTNAME (default: Admin)
* ADMIN_LASTNAME (default: Nanocloud)
* ADMIN_MAIL (default: admin@nanocloud.com)
* ADMIN_PASSWORD (default: admin)
* BACKEND_PORT (default: 8080)
* DATABASE_URI (mandatory)
* EXECUTION_SERVERS (mandatory)
* FRONT_DIR (mandatory)
* IAAS (default: qemu)
* LDAP_OU (default: OU=NanocloudUsers,DC=intra,DC=localdomain,DC=com)
* LDAP_PASSWORD (default: Nanocloud123+)
* LDAP_SERVER_URI (default: ldaps://Administrator:Nanocloud123+@iaas-module:6360)
* LDAP_USERNAME (default: CN=Administrator,CN=Users,DC=intra,DC=localdomain,DC=com)
* PLAZA_ADDRESS (default: iaas-module)
* PLAZA_PORT (default: 9090)
* PLAZA_USER_DIR (default: "C:\Users\%s\Desktop\Nanocloud")
* RDP_PORT (default: 3389)
* TRUST_PROXY (default: true)
* WINDOWS_DOMAIN (mandatory)
* WINDOWS_PASSWORD (mandatory)
* WINDOWS_USER (mandatory)

## Tests

To run backend unit tests:

````
docker-compose -f modules/docker-compose-build.yml run --rm nanocloud-backend make tests
````

To run frontend unit tests:

````
docker-compose -f modules/docker-compose-build.yml run --rm nanocloud-frontend ember test
````

To run API tests:

You need a fresh installation of Nanocloud community with a Windows VM running and ready to accept applications publish.

````
docker build -t nanocloud/testapi tests/api/
env NANOCLOUD-URL="localhost" docker run --net=host -e NANOCLOUD_HOST="${NANOCLOUR_URL}" --rm nanocloud/testapi
````

Replace localhost with the nanocloud's API host

## Roadmap

In future releases, we plan to add :

* Assign permission per application and per users.
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
