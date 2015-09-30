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

/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export var app = angular.module("hapticFrontend", ["flow", "ngRoute", "ngMaterial", "ngMdIcons", "ui.grid"]);

	app.config(["$routeProvider", "$mdThemingProvider", function($routeProvider: angular.route.IRouteProvider, $mdThemingProvider: angular.material.IThemingProvider) {

		$routeProvider
			.when("/", { templateUrl: "views/services.html", controller: "ServicesCtrl as servicesCtrl" })
			.when("/applications", { templateUrl: "views/applications.html", controller: "ApplicationsCtrl as applicationsCtrl" })
			.when("/users", { templateUrl: "views/users.html", controller: "UsersCtrl as usersCtrl" })
			.when("/stats", { templateUrl: "views/stats.html", controller: "StatsCtrl as statsCtrl" })
			.when("/capacity_planning", { templateUrl: "views/capacityPlanning.html", controller: "CapacityPlanningCtrl as capacityPlanningCtrl" })

			.otherwise({ redirectTo: "/" });

		$mdThemingProvider.theme("default")
			.primaryPalette("blue");

	}]);
}
