# Nanocloud dockerfiles

## Introduction

This project aims to deploy *Nanocloud Community* using an architecture split
into connected micro services.

These services are:

* a streaming server,
* a streaming client,
* a backend server,
* a frontend web server.

In this architecture, each service runs in a Docker container. The following
commands show how to deploy and run these containers.


## Configuration

The first step consists in downloading Nanocloud dockerfiles:

```
$ git clone https://github.com/Nanocloud/community.git
$ cd community/dockerfiles
```

And check content of all files in
*community/dockerfiles/nanocloud-backend/conf*.

### Web server
If you need to set a specific hostname for the **nginx** proxy, you can do it
by setting the directives *server_name* in the file **nginx/conf/nginx.conf**.


## Deployment

You can then deploy the services, in the *dockerfiles* directory, with the
following commands.

```
$ mkdir repos_guacamole; cd repos_guacamole/
$ git clone --depth 1 https://github.com/Nanocloud/noauth-logged.git
$ cd ..
```

```
$ mkdir repos_nanocloud; cd repos_nanocloud/
$ git clone --depth 1 https://github.com/Nanocloud/nanocloud.git
$ git clone --depth 1 https://github.com/Nanocloud/users.git
$ git clone --depth 1 https://github.com/Nanocloud/iaas.git
$ git clone --depth 1 https://github.com/Nanocloud/ldap.git
$ git clone --depth 1 https://github.com/Nanocloud/history.git
$ git clone --depth 1 https://github.com/Nanocloud/apps.git
$ cd ..
```

## Running

You need to have a working installation of **docker-compose** to build and
launch all the services.

Detailed information on **docker-compose** installation may be found
[here]("https://docs.docker.com/compose/install/").

Otherwise, you can just type the two following commands :

```
$ curl -L https://github.com/docker/compose/releases/download/1.4.2/docker-compose-`uname -s`-`uname -m` > docker-compose
$ chmod +x docker-compose
```

Then, start **Nanocloud Community** with :

```
$ ./docker-compose up -d
```

This command will build all docker containers, if they're not already, and start
them.

#### Licence

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
