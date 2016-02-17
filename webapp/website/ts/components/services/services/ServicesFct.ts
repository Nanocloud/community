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
/// <amd-dependency path="./ServicesSvc" />
import { IService, ServicesSvc } from "ServicesSvc";

"use strict";

export class ServicesFct {

	services: IService[] = [];
	private lock = false;
	private location: angular.ILocationService;

	static $inject = [
		"ServicesSvc",
		"$q",
		"$interval",
		"$location"
	];
	constructor(
		private servicesSvc: ServicesSvc,
		private $q: angular.IQService,
		$interval: ng.IIntervalService,
		$location: angular.ILocationService
	) {
		this.location = $location;
		this.loadServices();
		$interval(this.loadServices.bind(this), 10 * 1000);
	}

	loadServices(): angular.IPromise<void> {
		if (this.location.path() === "/admin/services") {
			if (this.lock) {
				return this.$q.resolve();
			}
			this.lock = true;
			return this.servicesSvc.getAll().then((services: IService[]) => {
				this.services = services;
				this.lock = false;
			});
		} else {
			return this.$q.when();
		}
	}

	getWindowsStatus(): angular.IPromise<boolean> {
		return this.servicesSvc.getAll().then((services: IService[]) => {
			let state = false;
			services.forEach(function(service: IService) {
				if (service.Ico === "windows" && service.Status === "running") {
					state = true;
				}
			});
			return state;
		});
	}
}

angular.module("haptic.services").service("ServicesFct", ServicesFct);
