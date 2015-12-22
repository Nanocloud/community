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
/// <amd-dependency path="../../applications/services/ApplicationsSvc" />
import { ApplicationsSvc, IApplication } from "../../applications/services/ApplicationsSvc";

"use strict";

export class PresenterCtrl {

	applications: any;
	user: string;

	static $inject = [
		"$state",
		"$cookies",
		"ApplicationsSvc"
	];

	constructor(
		private $state: angular.ui.IStateService,
		private $cookies: angular.cookies.ICookiesService,
		private appsSvc: ApplicationsSvc
	) {
		this.loadApplications();
		this.user = sessionStorage.getItem("user");
	}

	loadApplications(): angular.IPromise<void> {
		return this.appsSvc.getApplicationForUser().then((applications: IApplication[]) => {
			this.applications = applications;
		});
	}

	openApplication(application: IApplication, e: MouseEvent) {
		this.$cookies.remove("JSESSIONID");
		let appToken = btoa(application.ConnectionName + "\0c\0noauthlogged");
		let url = "/guacamole/#/client/" + appToken;
		if (localStorage["accessToken"]) {
			url += "?access_token=" + localStorage["accessToken"];
		}
		window.open(url, "_blank");
	}

	navigateTo(loc: string, e: MouseEvent) {
		window.open(loc, "_blank");
	}

	logout() {
		this.$state.go("logout");
	}

}

angular.module("haptic.presenter").controller("PresenterCtrl", PresenterCtrl);
