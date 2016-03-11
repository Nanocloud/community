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
/// <amd-dependency path="../../applications/controllers/DesktopCtrl" />
/// <amd-dependency path="../../applications/services/ApplicationsSvc" />
import { ApplicationsSvc, IApplication } from "../../applications/services/ApplicationsSvc";

"use strict";

export class PresenterCtrl {

	applications: any;
	user: string;
	private accessToken: string;

	static $inject = [
		"$state",
		"ApplicationsSvc",
		"$mdDialog",
		"$sce"
	];

	constructor(
		private $state: angular.ui.IStateService,
		private appsSvc: ApplicationsSvc,
		private $mdDialog: angular.material.IDialogService,
		private $sce: angular.ISCEService) {
		this.loadApplications();
		this.user = localStorage.getItem("user");
		this.accessToken = localStorage["accessToken"];
	}

	loadApplications(): angular.IPromise<void> {
		return this.appsSvc.getApplicationForUser().then((applications: IApplication[]) => {
			this.applications = applications;
		});
	}

	getAppUrl(application: IApplication) {
		let url = "/canva/#/canva/" + this.accessToken + "/" + application.alias;
		return this.$sce.trustAsResourceUrl(url);
	}

	openApplication(application: IApplication, e: MouseEvent) {
		let alias = "";
		if (application.alias === "Desktop") {
			alias = "hapticDesktop";
		} else {
			alias = application.alias;
		}
		let url = "/canva/#/canva/" + localStorage["accessToken"] + "/" + alias;
		window.open(url, "_blank");
	}

	navigateTo(loc: string, e: MouseEvent) {
		window.open(loc, "_blank");
	}

	logout() {
		this.$state.go("logout");
	}
	startDesktop(e: MouseEvent, url: string): angular.IPromise<any> {
		let o = this.getDefaultDesktopDlgOpt(e);
		o.locals = { desktop: this.$mdDialog, url: url };
		o.escapeToClose = false;
		o.onComplete = function() {
			$("#VDI")[0]["contentWindow"].focus();
		};
		return this.$mdDialog
			.show(o);
	}

	private getDefaultDesktopDlgOpt(e: MouseEvent): angular.material.IDialogOptions {
		return {
			controller: "DesktopCtrl",
			controllerAs: "desktopCtrl",
			templateUrl: "./js/components/applications/views/desktop.html",
			parent: angular.element(document.body),
			targetEvent: e
		};
	}
}

angular.module("haptic.presenter").controller("PresenterCtrl", PresenterCtrl);
