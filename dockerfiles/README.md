# Nanocloud dockerfiles

## Introduction
This project aims to deploy *Nanocloud Community* using an architecture split into connected micro services.

These services are:

- a streaming server,
- a streaming client,
- a backend server,
- a frontend web server.

In this architecture, each service runs in a Docker container. The following commands show how to deploy and run these containers.


## Configuration

The first step consists in downloading Nanocloud dockerfiles:

```
$ git clone ssh://git@git.nanocloud.com:7999/nanocloud/dockerfiles.git
$ cd dockerfiles
$ ls -F
docker-compose.yml  guacamole-client/  nanocloud-backend/  nginx/  README.md
```

### Execution server network
Before deploying the **nanocloud-backend** container, you must specify the IP adress of your **execution environment** instances. 
To do so, modify values of *AppServer.Server* and *AppServer.ExecutionServer* in the following lines of the file **nanocloud-backend/conf/config.json**.

```
...
    "AppServer": {
      "User" : "Administrator",
      "Server" : "<execution-manager-IPadress>",
      "ExecutionServers" : [
        "<execution-environment-IPadress>"
      ],
...
```

### Web server
If you need to set a specific hostname for the **nginx** proxy, you can do it by setting the directives *server_name* in the file **nginx/conf/nginx.conf**. 


## Deployment

You can then deploy the services, in the *dockerfiles* directory, with the following commands.

```
$ mkdir repos_guacamole; cd repos_guacamole/
$ git clone ssh://git@git.nanocloud.com:7999/nanocloud/noauth-logged.git
$ cd ..
```

```
$ mkdir repos_nanocloud; cd repos_nanocloud/
$ git clone ssh://git@git.nanocloud.com:7999/nanocloud/core.git
$ git clone ssh://git@git.nanocloud.com:7999/nanocloud/front.git
$ git clone ssh://git@git.nanocloud.com:7999/nanocloud/plugin_iaas.git
$ git clone ssh://git@git.nanocloud.com:7999/nanocloud/plugin_ldap.git
$ git clone ssh://git@git.nanocloud.com:7999/nanocloud/plugin_history.git
$ git clone ssh://git@git.nanocloud.com:7999/nanocloud/plugin_owncloud.git
$ cd ..
```


## Running

You need to have a working installation of **docker-compose** to build and launch all the services.
Detailed information on **docker-compose** installation may be found <a href="https://docs.docker.com/compose/install/" target="_blank">here</a>.

Otherwise, you can just type the two following commands :

```
$ curl -L https://github.com/docker/compose/releases/download/1.4.2/docker-compose-`uname -s`-`uname -m` > docker-compose
$ chmod +x docker-compose
```

Then, start **Nanocloud Community** with :

```
$ ./docker-compose up
```

This command will build all docker containers, if they're not already, and start them. 
Use CTRL-C in this command console to stop all the containers.



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
