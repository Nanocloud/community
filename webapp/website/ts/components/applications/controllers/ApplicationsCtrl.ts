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
/// <amd-dependency path="../services/ApplicationsSvc" />
/// <amd-dependency path="../../services/services/ServicesFct" />
/// <amd-dependency path="./ApplicationCtrl" />
/// <amd-dependency path="./DesktopCtrl" />
import { ApplicationsSvc, IApplication } from "../services/ApplicationsSvc";
import { ServicesFct } from "../../services/services/ServicesFct";

"use strict";

export class ApplicationsCtrl {

	applications: IApplication[];
	windowsState: boolean = false;
	loadWindowHasFinished: boolean = false;

	private accessToken: string;

	static $inject = [
		"ApplicationsSvc",
		"ServicesFct",
		"$mdDialog",
		"$sce"
	];

	constructor(
		private applicationsSrv: ApplicationsSvc,
		private servicesFct: ServicesFct,
		private $mdDialog: angular.material.IDialogService,
		private $sce: angular.ISCEService) {
		this.loadWindowHasFinished = false;
		this.servicesFct.getWindowsStatus().then((windowsState: boolean) => {
			this.loadWindowHasFinished = true;
			this.windowsState = windowsState;
		});
		this.accessToken = localStorage["accessToken"];
		this.applications = [];
		this.loadApplications();
	}

	loadApplications(): angular.IPromise<void> {
		return this.applicationsSrv.getAll().then((applications: IApplication[]) => {
			applications.forEach(function(application: IApplication) {
				if (application.alias !== "Desktop" && application.alias !== "hapticPowershell") {
					this.applications.push(application);
				}
			}.bind(this));
		});
	}

	startUnpublishApplication(e: MouseEvent, application: IApplication) {
		let o = this.$mdDialog.confirm()
			.parent(angular.element(document.body))
			.title("Unpublish application")
			.textContent("Are you sure you want to unpublish this application?")
			.ok("Yes")
			.cancel("No")
			.targetEvent(e);
		this.$mdDialog
			.show(o)
			.then(this.unpublishApplication.bind(this, application));
	}

	startRenameApplication(e: MouseEvent, application: IApplication) {
		let o = this.getDefaultRenameDlgOpt(e);
		o.locals = { app: application };
		return this.$mdDialog.show(o);
	}

	private getDefaultRenameDlgOpt(e: MouseEvent): angular.material.IDialogOptions {
		return {
			controller: "ApplicationCtrl",
			controllerAs: "applicationCtrl",
			templateUrl: "./js/components/applications/views/applicationrename.html",
			parent: angular.element(document.body),
			targetEvent: e
		};
	}

	unpublishApplication(application: IApplication) {
		this.applicationsSrv.unpublish(application);

		let i = _.findIndex(this.applications, (x: IApplication) => x.alias === application.alias);
		if (i >= 0) {
			this.applications.splice(i, 1);
		}
	}

	openApplication(e: MouseEvent, application: IApplication) {
		let appToken = btoa(application.alias + "\0c\0noauthlogged");
		let url = "/guacamole/#/client/" + appToken;
		if (localStorage["accessToken"]) {
			url += "?access_token=" + localStorage["accessToken"];
		}
		window.open(url, "_blank");
	}

	getDesktopUrl() {
		return this.$sce.trustAsResourceUrl("/canva/#/canva/" + this.accessToken + "/hapticDesktop");
	}

	getPublishUrl() {
		return this.$sce.trustAsResourceUrl("/canva/#/canva/" + this.accessToken + "/hapticPowershell");
	}

	getAppUrl(application: IApplication) {
		let url = "/canva/#/canva/" + this.accessToken + "/" + application.alias;
		return this.$sce.trustAsResourceUrl(url);
	}

	percentDone(file: any) {
		return Math.round(file._prevUploadedSize / file.size * 100).toString();
	}

	startDesktop(e: MouseEvent, url: string): angular.IPromise<any> {
		let o = this.getDefaultDesktopDlgOpt(e);
		o.locals = { desktop: this.$mdDialog, url: url };
		o.escapeToClose = false;
		o.onComplete = function() {
			setTimeout(function() {
				$("#VDI")[0]["contentWindow"].focus();
			}, 2000);
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

angular.module("haptic.applications").controller("ApplicationsCtrl", ApplicationsCtrl);
