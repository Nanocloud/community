/// <reference path='../../typings/tsd.d.ts' />

requirejs.config({
	baseUrl: "/js/",
	paths: {
		"jquery": "lib/jquery.min",
		"lodash": "lib/lodash.min",
		"angular": "lib/angular.min",
		"angular-cookies": "lib/angular-cookies.min",
		"angular-animate": "lib/angular-animate.min",
		"angular-aria": "lib/angular-aria.min",
		"angular-material": "lib/angular-material.min",
		"angular-material-icons": "lib/angular-material-icons.min",
		"angular-ui-grid": "lib/ui-grid.min",
		"angular-ui-route": "lib/angular-ui-router.min",
		"angular-ui-router-extras": "lib/ct-ui-router-extras.min",
		"ng-flow": "lib/ng-flow-standalone.min",
		
		"AmdTools": "components/core/services/AmdTools",
		"MainMenu": "components/core/services/MainMenu"
	},
	shim: {
		"angular": { exports: "angular", deps: ["jquery"] },
		"angular-route": { deps: ["angular"] },
		"angular-aria": { deps: ["angular"] },
		"angular-animate": { deps: ["angular"] },
		"angular-cookies": { deps: ["angular"] },
		"angular-material": { deps: ["angular", "angular-animate", "angular-aria"] },
		"angular-material-icons": { deps: ["angular-material"] },
		"angular-ui-grid": { deps: ["angular"] },
		"angular-ui-route": { deps: ["angular"] },
		"angular-ui-router-extras": { deps: ["angular-ui-route"] },
		"ng-flow": { deps: ["angular"] }
	},
	deps: ["haptic"],
	waitSeconds: 25
});
