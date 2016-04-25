#!/bin/bash

WIN_SERVER=$(cat windowsIP)
export WIN_SERVER

git clone https://github.com/Nanocloud/community.git

sed -i -e "s/=iaas-module/=$WIN_SERVER/" community/modules/docker-compose.yml
sed -i -e "s/IAAS=qemu/IAAS=manual/" community/modules/docker-compose.yml
sed -i -e "s/_PORT=1119/_PORT=22/" community/modules/docker-compose.yml
sed -i -e "s/iaas-module:6360/$WIN_SERVER:636/" community/modules/docker-compose.yml

n=0
until [ $n -ge 5 ]; do
	docker-compose -f community/modules/docker-compose-build.yml build
	docker-compose -f community/modules/docker-compose-build.yml up -d
	sleep 20
	CONTAINER_HOSTNAME=$(docker exec proxy hostname -I | awk '{print $1}') || "localhost"
	NANOCLOUD_STATUS=$(curl --output /dev/null --insecure --silent --write-out '%{http_code}\n' "https://${CONTAINER_HOSTNAME}")
	if [ "${NANOCLOUD_STATUS}" == "200" ]; then
		break
	fi
	docker-compose -f community/modules/docker-compose-build.yml down
	n=$((n+1))
done
