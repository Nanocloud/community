# Nanocloud community on Openstack

In order to avoid nested virtualization, a typical installation of Nanocloud on
openstack would require two virtual machines. One for Linux to host Nanocloud
and its data, and another for Windows to host applications.

This directory provides a script to automatically deploy such installation on
Openstack.
This installation script is used by Nanocloud's Bamboo to nightly deploy our
canary release in our effort to perform continuous integration.

## Requirements

*docker* is required to run the script. It will be the only software dependency.

You also need to have a valid Nanocloud Windows' image qcow2, either built
locally or downloaded from releases.nanocloud.org

## Deployment

Some parameters must be passed as environment variables.

- *DEPLOYMENT_OS_URL* openstack's URL
- *DEPLOYMENT_OS_USERNAME* user name to login
- *DEPLOYMENT_OS_PASSWORD* password to login
- *DEPLOYMENT_OS_PROJECT_ID* is the Openstack project identifier
- *DEPLOYMENT_OS_KEY_NAME* key name to start the Linux VM with

Optionnally these parameters can also be changed:

- *DEPLOYMENT_OS_INSTALL_SCRIPT_PATH* (default to ./installDocker.sh) script to be launch when Linux is available but before Windows is started
- *DEPLOYMENT_OS_WINDOWS_IMAGE_PATH* (default ./windows.qcow2 mounted by Docker) path to Nanocloud Windows image
- *DEPLOYMENT_OS_KEY_PATH* (default ./id_rsa mounted by Docker) private key that matched above key name to connect to Linux VM
- *DEPLOYMENT_OS_SSH_PORT* (default 22) SSH port on Linux Server

Then build the container:

````
docker build -t nanocloud/deployment-openstack .
````

Then run the container:

````
docker run -e DEPLOYMENT_OS_URL=http://my.openstack.com -e DEPLOYMENT_OS_PROJECT_ID=projectid -e DEPLOYMENT_OS_USERNAME=john -e DEPLOYMENT_OS_PASSWORD=pass -e DEPLOYMENT_OS_KEY_NAME=john -v /path/to/id_rsa:/opt/id_rsa -v /path/to/windows.qcow2:/opt/windows.qcow2 --rm nanocloud/deployment-openstack
````
