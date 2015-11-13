/// <reference path="../../../../typings/tsd.d.ts" />
/// <amd-dependency path="angular-ui-router-extras" />
import { overrideModuleRegisterer, registerCtrlFutureStates, getTemplateUrl } from "../core/services/AmdTools";
import { MainMenu } from "../core/services/MainMenu";

let componentName = "users";
let app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future"]);

let states: angular.ui.IState[] = [{
	name: "admin.users",
	url: "/users",
	controller: "UsersCtrl",
	controllerAs: "usersCtrl",
	templateUrl: getTemplateUrl(componentName, "users.html")
}];

MainMenu.add({
	stateName: states[0].name,
	title: "Users",
	ico: "people"
});

app.config(["$controllerProvider", "$provide", "$futureStateProvider", function(
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService,
	$futureStateProvider: any) {

	overrideModuleRegisterer(app, $controllerProvider, $provide);

	registerCtrlFutureStates(componentName, $futureStateProvider, states);

}]);
