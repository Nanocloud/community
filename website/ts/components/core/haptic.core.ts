/// <reference path="../../../../typings/tsd.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "./services/AmdTools";

let componentName = "core";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future"]);

app.config(["$controllerProvider", "$provide", "$futureStateProvider", "$urlRouterProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any,
	$urlRouterProvider: angular.ui.IUrlRouterProvider) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	$urlRouterProvider.otherwise("/");
	let states: angular.ui.IState[] = [{
		name: "admin",
		url: "/",
		controller: "MainCtrl",
		controllerAs: "mainCtrl",
		templateUrl: getTemplateUrl(componentName, "admin.html")
	}];
	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
