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
/// <amd-dependency path="../services/AuthenticationSvc" />
import { MainMenu, INavMenu } from "MainMenu";
import { AuthenticationSvc } from "../services/AuthenticationSvc";
import * as _ from "lodash";

"use strict";

export class MainCtrl {

	isLoading = false;
	user: string;
	menus: INavMenu[] = [];
	activeMenu: INavMenu = {};
	appInfo: any = {};

	static $inject = [
		"$state",
		"$mdSidenav",
		"AuthenticationSvc",
		"$rootScope",
		"$http"
	];
	constructor(
		private $state: angular.ui.IStateService,
		private $mdSidenav: angular.material.ISidenavService,
		private authSvc: AuthenticationSvc,
		$rootScope: angular.IRootScopeService,
		$http: angular.IHttpService
	) {
		this.user = sessionStorage.getItem("user");

		this.menus = _.sortBy(MainMenu.menus, (m: INavMenu) => m.title);
		this.checkMenuState($state.current);
		$rootScope.$on("$stateChangeSuccess", (event: angular.IAngularEvent, toState: angular.ui.IState) => {
			this.checkMenuState(toState);
		});
		
		$http.get("/api/version").success((res: any) => {
			this.appInfo = res;
		});
	}

	navigateTo(menu: INavMenu) {
		this.$mdSidenav("left").close();
		this.activeMenu = menu;
		this.$state.go(menu.stateName);
	}

	toggleMenu() {
		this.$mdSidenav("left").toggle();
	}

	checkMenuState(toState: angular.ui.IState) {
		let m = _.find(this.menus, (x: INavMenu) => x.stateName === toState.name);
		if (m) {
			this.activeMenu = m;
		}
	}

	logout() {
		this.authSvc.logout().then(() => {
			this.$state.go("login");
		});
	}

}

angular.module("haptic.core").controller("MainCtrl", MainCtrl);
