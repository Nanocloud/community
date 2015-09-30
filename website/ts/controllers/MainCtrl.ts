/*
 * Nanocloud community -- transform any application into SaaS solution
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

/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class MainCtrl {

		nav: models.MainNav;
		user: string;

		static $inject = [
			"$location"
		];
		constructor(
			private $location: angular.ILocationService
		) {
			this.nav = new models.MainNav();
			let m = _.find(this.nav.menus, (x: models.INavMenu) => x.url === $location.url());
			if (m) {
				this.nav.current = m;
			}
			this.user = sessionStorage.getItem("user");
		}

		navigateTo(menu: models.INavMenu) {
			this.$location.path(menu.url);
			this.nav.current = menu;
		}

	}

	app.controller("MainCtrl", MainCtrl);
}
