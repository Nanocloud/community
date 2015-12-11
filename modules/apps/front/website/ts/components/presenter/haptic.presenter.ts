/// <reference path="../../../../typings/tsd.d.ts" />
/// <reference path="../../core.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "AmdTools";

let componentName = "presenter";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future"]);

app.config(["$controllerProvider", "$provide", "$futureStateProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	let states: angular.ui.IState[] = [{
		name: "presenter",
		url: "/",
		controller: "PresenterCtrl",
		controllerAs: "presenterCtrl",
		templateUrl: getTemplateUrl(componentName, "presenter.html")
	}];
	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
