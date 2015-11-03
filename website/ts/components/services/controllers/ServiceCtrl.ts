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

/// <reference path='../../../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";
	
	class ServiceCtrl {

		service: IService;

		static $inject = [
			"ServicesService",
			"$mdDialog",
			"service"
		];
		constructor(
			private servicesSrv: ServicesService,
			private $mdDialog: angular.material.IDialogService,
			service: IService
		) {
			if (service) {
				this.service = angular.copy(service);
			}
		}

		accept(): void {
			if (this.servicesSrv.downloadStarted === false) {
				this.servicesSrv.download();
				this.$mdDialog.hide(this.service);
			} else {
				this.$mdDialog.cancel();
			}
		}

		close(): void {
			this.$mdDialog.cancel();
		}

	}

	angular.module("haptic.services").controller("ServiceCtrl", ServiceCtrl);
}
