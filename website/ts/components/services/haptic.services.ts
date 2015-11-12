/// <reference path="../../../../typings/tsd.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
/// <amd-dependency path="angular-ui-grid" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "../core/services/AmdTools";

let componentName = "services";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future", "ui.grid", "ui.grid.expandable"]);

app.config(["$controllerProvider", "$provide", "$futureStateProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	let states: angular.ui.IState[] = [{
		name: "admin.services",
		url: "services",
		controller: "ServicesCtrl",
		controllerAs: "serviceCtrl",
		templateUrl: getTemplateUrl(componentName, "services.html")
	}];
	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
