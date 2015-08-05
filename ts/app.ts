/// <reference path='../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export var app = angular.module("hapticFrontend", ["ngRoute", "ngMaterial", "ngMdIcons", "ui.grid"]);

	app.config(["$routeProvider", "$mdThemingProvider", function($routeProvider: angular.route.IRouteProvider, $mdThemingProvider: angular.material.IThemingProvider) {
		
		$routeProvider
			.when("/", { templateUrl: "views/home.html", controller: "HomeCtrl as homeCtrl" })
			.when("/login", { templateUrl: "views/login.html", controller: "LoginCtrl as loginCtrl" })
			.when("/users", { templateUrl: "views/users.html", controller: "UsersCtrl as usersCtrl" })
			.otherwise({ redirectTo: "/" });

		$mdThemingProvider.theme("default")
			.primaryPalette("blue");

	}]);
}
