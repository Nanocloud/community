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

/// <reference path="../../../../typings/tsd.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
/// <reference path="../../core.d.ts" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "AmdTools";

let componentName = "core";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future"]);

app.config(["$controllerProvider", "$provide", "$futureStateProvider", "$urlRouterProvider", "$urlMatcherFactoryProvider", "$httpProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any,
	$urlRouterProvider: angular.ui.IUrlRouterProvider,
	$urlMatcherFactoryProvider: angular.ui.IUrlMatcherFactory,
	$httpProvider: angular.IHttpProvider) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	// disable strict mode to be able to have "/foo/" == "/foo"
	$urlMatcherFactoryProvider.strictMode(false);
	
	// if a route wasn't found then go to the home
	$urlRouterProvider.otherwise(function($injector: angular.auto.IInjectorService, $location: angular.ILocationService): string {
		let prefix = "/admin";
		if ($location.url().slice(0, prefix.length) === prefix) {
			return prefix + "/services"; // default admin page
		} else {
			return "/"; // default normal page
		}
	});

	// admin parent state, all admin routes are passed by this
	let states: angular.ui.IState[] = [{
		abstract: true,
		name: "admin",
		url: "/admin",
		controller: "MainCtrl",
		controllerAs: "mainCtrl",
		templateUrl: getTemplateUrl(componentName, "admin.html")
	}];
	registerCtrlFutureStates(componentName, $futureStateProvider, states);

	// if an oauth authentication is found then set it in all http headers
	if (localStorage["accessToken"]) {
		$httpProvider.defaults.headers.common["Authorization"] = "Bearer " + localStorage["accessToken"];
	}

	// allows to have a global spinner for ajax requests
	$httpProvider.interceptors.push(function() {
		return {
			"request": function(config: any) {
				let spn = document.getElementById("coreSpinner");
				if (spn) {
					spn.style.visibility = "visible";
				}
				return config;
			},
			"response": function(response: any) {
				let spn = document.getElementById("coreSpinner");
				if (spn) {
					spn.style.visibility = "hidden";
				}
				return response;
			}
		};
	});

	// global error handler
	$httpProvider.interceptors.push(["$injector", "$q", function($injector: angular.auto.IInjectorService, $q: angular.IQService) {
		return {
			"responseError": function(rejection: angular.IHttpPromiseCallbackArg<any>) {
				if (rejection.status === 401 || rejection.status === 403) {
					let $location = <angular.ILocationService>$injector.get("$location");
					$location.path("/login");
				} else {
					let $mdToast = <angular.material.IToastService>$injector.get("$mdToast");
					$mdToast.show(
						$mdToast.simple()
							.textContent(rejection.statusText)
							.position("top right")
					);
				}
				return $q.reject(rejection);
			}
		};
	}]);

}]);
