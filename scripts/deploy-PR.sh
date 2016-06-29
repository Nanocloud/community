#!/bin/sh -e
#
# Nanocloud Community, a comprehensive platform to turn any application
# into a cloud solution.
#
# Copyright (C) 2015 Nanocloud Software
#
# This file is part of Nanocloud community.
#
# Nanocloud community is free software; you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as
# published by the Free Software Foundation, either version 3 of the
# License, or (at your option) any later version.
#
# Nanocloud community is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <http://www.gnu.org/licenses/>.

BRANCH=${1}

mkdir .deploy ; cd .deploy

wget http://releases.nanocloud.org:8080/builds/$BRANCH/nanocloud-frontend.tar
wget http://releases.nanocloud.org:8080/builds/$BRANCH/nanocloud-backend.tar

for file in *.tar; do
    docker load -i $file
done

cd ../ && rm -rf .deploy
