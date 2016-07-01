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

# $GITHUB_PASSWORD Github access token
# $1 = STATE "success" "error" or "failure
# $2 = $DESCRIPTION
# $3 = $CONTEXT
# $4 = $URL Url to access test result
# $5 = $REPO Github format. Example: Nanocloud/community
# $6 = $SHA Commit concerned

curl -H "Authorization: token ${GITHUB_PASSWORD}" --request POST -k --data "{\"state\": \"${1}\", \"description\": \"${2}\", \"context\": \"${3}\", \"target_url\": \"${4}\"}" https://api.github.com/repos/${5}/statuses/${6}
