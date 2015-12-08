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
/// <amd-dependency path="../services/ServicesSvc" />
/// <amd-dependency path="./ServiceCtrl" />
import { ServicesSvc, IService } from "../services/ServicesSvc";

"use strict";

export class ServicesCtrl {

	services: any;
	colors: any;

	static $inject = [
		"ServicesSvc",
		"$mdDialog",
		"$interval",
		"$scope"
	];

	constructor(
		private servicesSrv: ServicesSvc,
		private $mdDialog: angular.material.IDialogService,
		$interval: ng.IIntervalService,
		$scope: angular.IScope
	) {
		this.colors = {
			downloading: "#4183D7",
			available: "#A2DED0",
			booting: "#EB9532",
			running: "#26A65B"
		};

		this.loadServices();
		let intervalPrms = $interval(this.loadServices.bind(this), 10 * 1000);
		$scope.$on("$destroy", function () {
			$interval.cancel(intervalPrms);
		});
	}

	loadServices(): angular.IPromise<void> {
		return this.servicesSrv.getAll().then((services: IService[]) => {
			this.services = services;
		});
	}

	startWindowsDownload(e: MouseEvent, service: IService) {
		let o = this.getDefaultServiceDlgOpt(e);
		o.locals = { service: service };
		return this.$mdDialog.show(o);
	}

	toggle(service: IService) {
		if (! service.Locked) {
			if (service.Status === "running") {
				this.servicesSrv.startStopService(service);
			} else {
				this.servicesSrv.start(service);
			}
		}
	}

	private getDefaultServiceDlgOpt(e: MouseEvent): angular.material.IDialogOptions {
		return {
			controller: "ServiceCtrl",
			controllerAs: "serviceCtrl",
			templateUrl: "./js/components/services/views/service.html",
			parent: angular.element(document.body),
			targetEvent: e
		};
	}

}

angular.module("haptic.services").controller("ServicesCtrl", ServicesCtrl);
