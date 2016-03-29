
export WIN_SERVER=$(cat windowsIP)

git clone https://github.com/Nanocloud/community.git

sed -i -e "s/=iaas-module/=$WIN_SERVER/" community/modules/docker-compose.yml
sed -i -e "s/IAAS=qemu/IAAS=manual/" community/modules/docker-compose.yml
sed -i -e "s/_PORT=1119/_PORT=22/" community/modules/docker-compose.yml

docker-compose -f community/modules/docker-compose-build.yml up -d
