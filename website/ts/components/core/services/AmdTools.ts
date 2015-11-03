/// <reference path="../../../../../typings/tsd.d.ts" />

// overloads that allows to register dynamically after module.config
export function overrideModuleRegisterer(
	app: angular.IModule,
	$controllerProvider: angular.IControllerProvider,
	$provide: angular.auto.IProvideService): void {
	// override controller
	(<any>app).controller = function(name: string, controllerConstructor: Function): angular.IModule {
		$controllerProvider.register(name, controllerConstructor);
		return app;
	};
	// override service
	(<any>app).service = function(name: string, serviceConstructor: Function): angular.IModule {
		$provide.service(name, serviceConstructor);
		return app;
	};
}

// allows to load the controller before navigating to the router
export var requireCtrlStateFactory = ["$q", "futureState", function($q: angular.IQService, futureState: angular.ui.IState) {
	let defer = $q.defer();
	requirejs(["components/core/controllers/" + futureState.controller], function() {
		defer.resolve(futureState);
	});
	return defer.promise;
}];

// register states to placeholders that will be replaced with a full UI-Router state when navigated to
export function registerCtrlFutureStates($futureStateProvider: any, states: angular.ui.IState[]): void {
	$futureStateProvider.stateFactory("requireCtrl", requireCtrlStateFactory);
	for (var state of states) {
		(<any>state).type = "requireCtrl";
		state.templateUrl = "./js/components/core/views/" + state.templateUrl;
		$futureStateProvider.futureState(state);
	}
}