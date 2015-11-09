/// <reference path="../../../../typings/tsd.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "../core/services/AmdTools";

let componentName = "users";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future"]);

app.config(["$controllerProvider", "$provide", "$futureStateProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	let states: angular.ui.IState[] = [{
		name: "admin.users",
		url: "users",
		controller: "UsersCtrl",
		controllerAs: "usersCtrl",
		templateUrl: getTemplateUrl(componentName, "users.html")
	}];
	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
