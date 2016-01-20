define(["require", "exports", "AmdTools", "angular-ui-router-extras"], function (require, exports, AmdTools_1) {
    var componentName = "presenter";
    var app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future"]);
    app.config(["$controllerProvider", "$provide", "$futureStateProvider", function ($controllerProvider, $provide, $futureStateProvider) {
            AmdTools_1.overrideModuleRegisterer(app, $controllerProvider, $provide);
            var states = [{
                    name: "presenter",
                    url: "/",
                    controller: "PresenterCtrl",
                    controllerAs: "presenterCtrl",
                    templateUrl: AmdTools_1.getTemplateUrl(componentName, "presenter.html")
                }];
            AmdTools_1.registerCtrlFutureStates(componentName, $futureStateProvider, states);
        }]);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaGFwdGljLnByZXNlbnRlci5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uL3RzL2NvbXBvbmVudHMvcHJlc2VudGVyL2hhcHRpYy5wcmVzZW50ZXIudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtJQTJCQSxJQUFJLGFBQWEsR0FBRyxXQUFXLENBQUM7SUFDaEMsSUFBSSxHQUFHLEdBQUcsT0FBTyxDQUFDLE1BQU0sQ0FBQyxTQUFTLEdBQUcsYUFBYSxFQUFFLENBQUMsNEJBQTRCLENBQUMsQ0FBQyxDQUFDO0lBRXBGLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQyxxQkFBcUIsRUFBRSxVQUFVLEVBQUUsc0JBQXNCLEVBQUUsVUFDdEUsbUJBQWdELEVBQ2hELFFBQXNDLEVBQ3RDLG9CQUF5QjtZQUV6QixtQ0FBd0IsQ0FBQyxHQUFHLEVBQUUsbUJBQW1CLEVBQUUsUUFBUSxDQUFDLENBQUM7WUFFN0QsSUFBSSxNQUFNLEdBQXdCLENBQUM7b0JBQ2xDLElBQUksRUFBRSxXQUFXO29CQUNqQixHQUFHLEVBQUUsR0FBRztvQkFDUixVQUFVLEVBQUUsZUFBZTtvQkFDM0IsWUFBWSxFQUFFLGVBQWU7b0JBQzdCLFdBQVcsRUFBRSx5QkFBYyxDQUFDLGFBQWEsRUFBRSxnQkFBZ0IsQ0FBQztpQkFDNUQsQ0FBQyxDQUFDO1lBQ0gsbUNBQXdCLENBQUMsYUFBYSxFQUFFLG9CQUFvQixFQUFFLE1BQU0sQ0FBQyxDQUFDO1FBRXZFLENBQUMsQ0FBQyxDQUFDLENBQUMifQ==