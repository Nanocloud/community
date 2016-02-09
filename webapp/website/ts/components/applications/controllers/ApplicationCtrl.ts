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
import { ApplicationsSvc, IApplication } from "../services/ApplicationsSvc";

"use strict";

export class ApplicationCtrl {

	app: IApplication;
	name: string;

	static $inject = [
		"ApplicationsSvc",
		"$mdDialog",
		"app",
		"$state",
	];

	constructor(
		private applicationsSvc: ApplicationsSvc,
		private $mdDialog: angular.material.IDialogService,
		app: IApplication,
		private $state: ng.ui.IStateService
	) {
		this.app = angular.copy(app);
		this.name = this.app.display_name;
	}

	resetName(): void {
		this.name = this.app.alias;
	}

	save(): void {
		this.applicationsSvc.changeName(this.app, this.name).then(function() {
			this.$mdDialog.cancel();
			this.$state.reload();
		}.bind(this));
	}
}

angular.module("haptic.applications").controller("ApplicationCtrl", ApplicationCtrl);
