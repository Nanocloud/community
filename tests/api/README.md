# Nanocloud community functional tests

This directory provides API functional tests.

## Requirements

You need a fresh installation of Nanocloud community with a Windows VM running
and ready to accept applications.

## Launch Tests

```
docker build -t nanocloud/testapi .
```

Then run the container:

```
NANOCLOUD_URL="localhost"
docker run --net=host -e NANOCLOUD_HOST="${NANOCLOUR_URL}" nanocloud/testapi
```
