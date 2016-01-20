define(["require", "exports", "AmdTools", "angular-ui-router-extras"], function (require, exports, AmdTools_1) {
    var componentName = "login";
    var app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future"]);
    app.config(["$controllerProvider", "$provide", "$futureStateProvider", function ($controllerProvider, $provide, $futureStateProvider) {
            AmdTools_1.overrideModuleRegisterer(app, $controllerProvider, $provide);
            var states = [{
                    name: "login",
                    url: "/login",
                    controller: "LoginCtrl",
                    controllerAs: "loginCtrl",
                    templateUrl: AmdTools_1.getTemplateUrl(componentName, "login.html")
                }, {
                    name: "logout",
                    url: "/logout",
                    controller: "LoginCtrl",
                    controllerAs: "loginCtrl",
                    templateUrl: AmdTools_1.getTemplateUrl(componentName, "login.html"),
                    params: { logout: true }
                }];
            AmdTools_1.registerCtrlFutureStates(componentName, $futureStateProvider, states);
        }]);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaGFwdGljLmxvZ2luLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vdHMvY29tcG9uZW50cy9sb2dpbi9oYXB0aWMubG9naW4udHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtJQTJCQSxJQUFJLGFBQWEsR0FBRyxPQUFPLENBQUM7SUFDNUIsSUFBSSxHQUFHLEdBQUcsT0FBTyxDQUFDLE1BQU0sQ0FBQyxTQUFTLEdBQUcsYUFBYSxFQUFFLENBQUMsNEJBQTRCLENBQUMsQ0FBQyxDQUFDO0lBRXBGLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQyxxQkFBcUIsRUFBRSxVQUFVLEVBQUUsc0JBQXNCLEVBQUUsVUFDdEUsbUJBQWdELEVBQ2hELFFBQXNDLEVBQ3RDLG9CQUF5QjtZQUV6QixtQ0FBd0IsQ0FBQyxHQUFHLEVBQUUsbUJBQW1CLEVBQUUsUUFBUSxDQUFDLENBQUM7WUFFN0QsSUFBSSxNQUFNLEdBQXdCLENBQUM7b0JBQ2xDLElBQUksRUFBRSxPQUFPO29CQUNiLEdBQUcsRUFBRSxRQUFRO29CQUNiLFVBQVUsRUFBRSxXQUFXO29CQUN2QixZQUFZLEVBQUUsV0FBVztvQkFDekIsV0FBVyxFQUFFLHlCQUFjLENBQUMsYUFBYSxFQUFFLFlBQVksQ0FBQztpQkFDeEQsRUFBRTtvQkFDRixJQUFJLEVBQUUsUUFBUTtvQkFDZCxHQUFHLEVBQUUsU0FBUztvQkFDZCxVQUFVLEVBQUUsV0FBVztvQkFDdkIsWUFBWSxFQUFFLFdBQVc7b0JBQ3pCLFdBQVcsRUFBRSx5QkFBYyxDQUFDLGFBQWEsRUFBRSxZQUFZLENBQUM7b0JBQ3hELE1BQU0sRUFBRSxFQUFFLE1BQU0sRUFBRSxJQUFJLEVBQUU7aUJBQ3hCLENBQUMsQ0FBQztZQUNILG1DQUF3QixDQUFDLGFBQWEsRUFBRSxvQkFBb0IsRUFBRSxNQUFNLENBQUMsQ0FBQztRQUV2RSxDQUFDLENBQUMsQ0FBQyxDQUFDIn0=