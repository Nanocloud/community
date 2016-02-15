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
        "angular-ui-route": { deps: ["angular"] },
        "angular-ui-router-extras": { deps: ["angular-ui-route"] },
        "ng-flow": { deps: ["angular"] }
    },
    deps: ["haptic"],
    waitSeconds: 25
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiYm9vdHN0cmFwLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vdHMvYm9vdHN0cmFwLnRzIl0sIm5hbWVzIjpbXSwibWFwcGluZ3MiOiJBQXdCQSxTQUFTLENBQUMsTUFBTSxDQUFDO0lBQ2hCLE9BQU8sRUFBRSxNQUFNO0lBQ2YsS0FBSyxFQUFFO1FBQ04sUUFBUSxFQUFFLGdCQUFnQjtRQUMxQixRQUFRLEVBQUUsZ0JBQWdCO1FBQzFCLFNBQVMsRUFBRSxpQkFBaUI7UUFDNUIsaUJBQWlCLEVBQUUseUJBQXlCO1FBQzVDLGlCQUFpQixFQUFFLHlCQUF5QjtRQUM1QyxjQUFjLEVBQUUsc0JBQXNCO1FBQ3RDLGtCQUFrQixFQUFFLDBCQUEwQjtRQUM5Qyx3QkFBd0IsRUFBRSxnQ0FBZ0M7UUFDMUQsaUJBQWlCLEVBQUUsaUJBQWlCO1FBQ3BDLGtCQUFrQixFQUFFLDJCQUEyQjtRQUMvQywwQkFBMEIsRUFBRSw2QkFBNkI7UUFDekQsU0FBUyxFQUFFLDRCQUE0QjtRQUV2QyxVQUFVLEVBQUUsbUNBQW1DO1FBQy9DLFVBQVUsRUFBRSxtQ0FBbUM7S0FDL0M7SUFDRCxJQUFJLEVBQUU7UUFDTCxTQUFTLEVBQUUsRUFBRSxPQUFPLEVBQUUsU0FBUyxFQUFFLElBQUksRUFBRSxDQUFDLFFBQVEsQ0FBQyxFQUFFO1FBQ25ELGVBQWUsRUFBRSxFQUFFLElBQUksRUFBRSxDQUFDLFNBQVMsQ0FBQyxFQUFFO1FBQ3RDLGNBQWMsRUFBRSxFQUFFLElBQUksRUFBRSxDQUFDLFNBQVMsQ0FBQyxFQUFFO1FBQ3JDLGlCQUFpQixFQUFFLEVBQUUsSUFBSSxFQUFFLENBQUMsU0FBUyxDQUFDLEVBQUU7UUFDeEMsaUJBQWlCLEVBQUUsRUFBRSxJQUFJLEVBQUUsQ0FBQyxTQUFTLENBQUMsRUFBRTtRQUN4QyxrQkFBa0IsRUFBRSxFQUFFLElBQUksRUFBRSxDQUFDLFNBQVMsRUFBRSxpQkFBaUIsRUFBRSxjQUFjLENBQUMsRUFBRTtRQUM1RSx3QkFBd0IsRUFBRSxFQUFFLElBQUksRUFBRSxDQUFDLGtCQUFrQixDQUFDLEVBQUU7UUFDeEQsaUJBQWlCLEVBQUUsRUFBRSxJQUFJLEVBQUUsQ0FBQyxTQUFTLENBQUMsRUFBRTtRQUN4QyxrQkFBa0IsRUFBRSxFQUFFLElBQUksRUFBRSxDQUFDLFNBQVMsQ0FBQyxFQUFFO1FBQ3pDLDBCQUEwQixFQUFFLEVBQUUsSUFBSSxFQUFFLENBQUMsa0JBQWtCLENBQUMsRUFBRTtRQUMxRCxTQUFTLEVBQUUsRUFBRSxJQUFJLEVBQUUsQ0FBQyxTQUFTLENBQUMsRUFBRTtLQUNoQztJQUNELElBQUksRUFBRSxDQUFDLFFBQVEsQ0FBQztJQUNoQixXQUFXLEVBQUUsRUFBRTtDQUNmLENBQUMsQ0FBQyJ9