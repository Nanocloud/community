define(["require", "exports", "AmdTools", "MainMenu", "angular-ui-router-extras", "angular-cookies", "angular-ui-grid", "ng-flow"], function (require, exports, AmdTools_1, MainMenu_1) {
    var componentName = "applications";
    var app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future", "ui.grid", "ngCookies", "flow"]);
    var states = [{
            name: "admin.applications",
            url: "/applications",
            controller: "ApplicationsCtrl",
            controllerAs: "applicationsCtrl",
            templateUrl: AmdTools_1.getTemplateUrl(componentName, "applications.html")
        }];
    MainMenu_1.MainMenu.add({
        stateName: states[0].name,
        title: "Applications",
        ico: "apps"
    });
    app.config(["$controllerProvider", "$provide", "$futureStateProvider", "flowFactoryProvider", function ($controllerProvider, $provide, $futureStateProvider, flowFactoryProvider) {
            AmdTools_1.overrideModuleRegisterer(app, $controllerProvider, $provide);
            AmdTools_1.registerCtrlFutureStates(componentName, $futureStateProvider, states);
            flowFactoryProvider.defaults = {
                headers: {
                    "Authorization": "Bearer " + localStorage["accessToken"]
                }
            };
        }]);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaGFwdGljLmFwcGxpY2F0aW9ucy5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3RzL2NvbXBvbmVudHMvYXBwbGljYXRpb25zL2hhcHRpYy5hcHBsaWNhdGlvbnMudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtJQStCQSxJQUFJLGFBQWEsR0FBRyxjQUFjLENBQUM7SUFDbkMsSUFBSSxHQUFHLEdBQUcsT0FBTyxDQUFDLE1BQU0sQ0FBQyxTQUFTLEdBQUcsYUFBYSxFQUFFLENBQUMsNEJBQTRCLEVBQUUsU0FBUyxFQUFFLFdBQVcsRUFBRSxNQUFNLENBQUMsQ0FBQyxDQUFDO0lBRXBILElBQUksTUFBTSxHQUF3QixDQUFDO1lBQ2xDLElBQUksRUFBRSxvQkFBb0I7WUFDMUIsR0FBRyxFQUFFLGVBQWU7WUFDcEIsVUFBVSxFQUFFLGtCQUFrQjtZQUM5QixZQUFZLEVBQUUsa0JBQWtCO1lBQ2hDLFdBQVcsRUFBRSx5QkFBYyxDQUFDLGFBQWEsRUFBRSxtQkFBbUIsQ0FBQztTQUMvRCxDQUFDLENBQUM7SUFFSCxtQkFBUSxDQUFDLEdBQUcsQ0FBQztRQUNaLFNBQVMsRUFBRSxNQUFNLENBQUMsQ0FBQyxDQUFDLENBQUMsSUFBSTtRQUN6QixLQUFLLEVBQUUsY0FBYztRQUNyQixHQUFHLEVBQUUsTUFBTTtLQUNYLENBQUMsQ0FBQztJQUVILEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQyxxQkFBcUIsRUFBRSxVQUFVLEVBQUUsc0JBQXNCLEVBQUUscUJBQXFCLEVBQUUsVUFDN0YsbUJBQWdELEVBQ2hELFFBQXNDLEVBQ3RDLG9CQUF5QixFQUN6QixtQkFBd0I7WUFFeEIsbUNBQXdCLENBQUMsR0FBRyxFQUFFLG1CQUFtQixFQUFFLFFBQVEsQ0FBQyxDQUFDO1lBRTdELG1DQUF3QixDQUFDLGFBQWEsRUFBRSxvQkFBb0IsRUFBRSxNQUFNLENBQUMsQ0FBQztZQUV0RSxtQkFBbUIsQ0FBQyxRQUFRLEdBQUc7Z0JBQzlCLE9BQU8sRUFBRTtvQkFDUixlQUFlLEVBQUUsU0FBUyxHQUFHLFlBQVksQ0FBQyxhQUFhLENBQUM7aUJBQ3hEO2FBQ0QsQ0FBQztRQUVILENBQUMsQ0FBQyxDQUFDLENBQUMifQ==