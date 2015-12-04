/// <reference path="../../../../typings/tsd.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "./services/AmdTools";

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

	$urlMatcherFactoryProvider.strictMode(false);
	$urlRouterProvider.otherwise(function($injector: angular.auto.IInjectorService, $location: angular.ILocationService): string {
		let prefix = "/admin";
		if ($location.url().slice(0, prefix.length) === prefix) {
			return prefix + "/services"; // default admin page
		} else {
			return "/";
		}
	});

	let states: angular.ui.IState[] = [{
		abstract: true,
		name: "admin",
		url: "/admin",
		controller: "MainCtrl",
		controllerAs: "mainCtrl",
		templateUrl: getTemplateUrl(componentName, "admin.html")
	}];
	registerCtrlFutureStates(componentName, $futureStateProvider, states);

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

}]);
