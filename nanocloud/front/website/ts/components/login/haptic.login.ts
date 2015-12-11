/// <reference path="../../../../typings/tsd.d.ts" />
/// <reference path="../../core.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "AmdTools";

let componentName = "login";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future"]);

app.config(["$controllerProvider", "$provide", "$futureStateProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	let states: angular.ui.IState[] = [{
		name: "login",
		url: "/login",
		controller: "LoginCtrl",
		controllerAs: "loginCtrl",
		templateUrl: getTemplateUrl(componentName, "login.html")
	}];
	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
