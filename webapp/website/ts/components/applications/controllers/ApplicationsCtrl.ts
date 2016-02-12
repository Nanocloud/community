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
/// <amd-dependency path="./ApplicationCtrl" />
import { ApplicationsSvc, IApplication } from "../services/ApplicationsSvc";

"use strict";

export class ApplicationsCtrl {

	gridOptions: any;

	static $inject = [
		"ApplicationsSvc",
		"$mdDialog"
	];

	constructor(
		private applicationsSrv: ApplicationsSvc,
		private $mdDialog: angular.material.IDialogService
	) {
		this.gridOptions = {
			data: [],
			rowHeight: 36,
			columnDefs: [
				{ field: "Application" },
				{
					name: "actions",
					displayName: "",
					enableColumnMenu: false,
					cellTemplate: `
						<md-button ng-click='grid.appScope.applicationsCtrl.openApplication($event, row.entity)'>
							<ng-md-icon icon='pageview' size='14'></ng-md-icon> Open
						</md-button>
						<md-button ng-click='grid.appScope.applicationsCtrl.startRenameApplication($event, row.entity)'>
							<ng-md-icon icon='edit' size='14'></ng-md-icon> Edit name
						</md-button>
						<md-button ng-click='grid.appScope.applicationsCtrl.startUnpublishApplication($event, row.entity)'>
							<ng-md-icon icon='delete' size='14'></ng-md-icon> Unpublish
						</md-button>`
				}
			]	
		};

		this.loadApplications();
	}

	get applications(): IApplication[] {
		return this.gridOptions.data;
	}
	set applications(value: IApplication[]) {
		this.gridOptions.data = value;
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
		let appToken = btoa(application.Application + "\0c\0noauthlogged");
		let url = "/guacamole/#/client/" + appToken;
		if (localStorage["accessToken"]) {
			url += "?access_token=" + localStorage["accessToken"];
		}
		window.open(url, "_blank");
	}

	percentDone(file: any) {
		return Math.round(file._prevUploadedSize / file.size * 100).toString();
	}

}

angular.module("haptic.applications").controller("ApplicationsCtrl", ApplicationsCtrl);
