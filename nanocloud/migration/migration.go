/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package migration

import (
	"github.com/Nanocloud/community/nanocloud/migration/apps"
	"github.com/Nanocloud/community/nanocloud/migration/config"
	"github.com/Nanocloud/community/nanocloud/migration/history"
	"github.com/Nanocloud/community/nanocloud/migration/machines"
	"github.com/Nanocloud/community/nanocloud/migration/oauth"
	"github.com/Nanocloud/community/nanocloud/migration/users"

	log "github.com/Sirupsen/logrus"
)

func Migrate() error {
	err := users.Migrate()
	if err != nil {
		log.Error("users migration failed")
		return err
	}

	err = oauth.Migrate()
	if err != nil {
		log.Error("oauth migration failed")
		return err
	}

	err = apps.Migrate()
	if err != nil {
		log.Error("apps migration failed")
		return err
	}

	err = history.Migrate()
	if err != nil {
		log.Error("history migration failed")
		return err
	}

	err = machines.Migrate()
	if err != nil {
		log.Error("machines migration failed")
		return err
	}

	err = config.Migrate()
	if err != nil {
		log.Error("config migration failed")
		return err
	}

	return nil
}
