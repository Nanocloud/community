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

/// <reference path="../../../../../typings/tsd.d.ts" />
import { MainMenu, INavMenu } from "../services/MainMenu";
import * as _ from "lodash";

"use strict";

export class MainCtrl {

	user: string;
	menus: INavMenu[] = [];
	activeMenu: INavMenu = {};

	static $inject = [
		"$state",
		"$mdSidenav",
		"$rootScope"
	];
	constructor(
		private $state: angular.ui.IStateService,
		private $mdSidenav: angular.material.ISidenavService,
		private $rootScope: angular.IRootScopeService
	) {
		this.menus = _.sortBy(MainMenu.menus, (m: INavMenu) => m.title);
		
		this.checkMenu(null, $state.current);
		$rootScope.$on("$stateChangeSuccess", this.checkMenu.bind(this));

		this.user = sessionStorage.getItem("user");
	}

	navigateTo(menu: INavMenu) {
		this.$mdSidenav("left").close();
		this.activeMenu = menu;
		this.$state.go(menu.stateName);
	}

	toggleMenu() {
		this.$mdSidenav("left").toggle();
	}

	checkMenu(event: angular.IAngularEvent, toState: angular.ui.IState) {
		let m = _.find(this.menus, (x: INavMenu) => x.stateName === toState.name);
		if (m) {
			this.activeMenu = m;
		}
	}

}

angular.module("haptic.core").controller("MainCtrl", MainCtrl);
