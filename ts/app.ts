/// <reference path='../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export var app = angular.module("hapticFrontend", ["ngRoute", "ngMaterial", "ngMdIcons", "ui.grid"]);

	app.config(["$routeProvider", "$mdThemingProvider", function($routeProvider: angular.route.IRouteProvider, $mdThemingProvider: angular.material.IThemingProvider) {

		$routeProvider
			.when("/", { templateUrl: "views/services.html", controller: "ServicesCtrl as servicesCtrl" })
			.when("/login", { templateUrl: "views/login.html", controller: "LoginCtrl as loginCtrl" })
			.when("/applications", { templateUrl: "views/applications.html", controller: "ApplicationsCtrl as applicationsCtrl" })
			.when("/users", { templateUrl: "views/users.html", controller: "UsersCtrl as usersCtrl" })
			.when("/stats", { templateUrl: "views/stats.html", controller: "StatsCtrl as statsCtrl" })
			.when("/capacity_planning", { templateUrl: "views/capacityPlanning.html", controller: "CapacityPlanningCtrl as capacityPlanningCtrl" })

			.otherwise({ redirectTo: "/" });

		$mdThemingProvider.theme("default")
			.primaryPalette("blue");

	}]);
}
