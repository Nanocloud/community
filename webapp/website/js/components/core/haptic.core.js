define(["require", "exports", "AmdTools", "angular-ui-router-extras"], function (require, exports, AmdTools_1) {
    var componentName = "core";
    var app = angular.module("haptic." + componentName, ["ct.ui.router.extras.future"]);
    app.config(["$controllerProvider", "$provide", "$futureStateProvider", "$urlRouterProvider", "$urlMatcherFactoryProvider", "$httpProvider", function ($controllerProvider, $provide, $futureStateProvider, $urlRouterProvider, $urlMatcherFactoryProvider, $httpProvider) {
            AmdTools_1.overrideModuleRegisterer(app, $controllerProvider, $provide);
            $urlMatcherFactoryProvider.strictMode(false);
            $urlRouterProvider.otherwise(function ($injector, $location) {
                var prefix = "/admin";
                if ($location.url().slice(0, prefix.length) === prefix) {
                    return prefix + "/services";
                }
                else {
                    return "/";
                }
            });
            var states = [{
                    abstract: true,
                    name: "admin",
                    url: "/admin",
                    controller: "MainCtrl",
                    controllerAs: "mainCtrl",
                    templateUrl: AmdTools_1.getTemplateUrl(componentName, "admin.html")
                }];
            AmdTools_1.registerCtrlFutureStates(componentName, $futureStateProvider, states);
            if (localStorage["accessToken"]) {
                $httpProvider.defaults.headers.common["Authorization"] = "Bearer " + localStorage["accessToken"];
            }
            $httpProvider.interceptors.push(function () {
                return {
                    "request": function (config) {
                        var spn = document.getElementById("coreSpinner");
                        if (spn) {
                            spn.style.visibility = "visible";
                        }
                        return config;
                    },
                    "response": function (response) {
                        var spn = document.getElementById("coreSpinner");
                        if (spn) {
                            spn.style.visibility = "hidden";
                        }
                        return response;
                    }
                };
            });
            $httpProvider.interceptors.push(["$injector", "$q", function ($injector, $q) {
                    return {
                        "responseError": function (rejection) {
                            if (rejection.status === 401 || rejection.status === 403) {
                                var $location = $injector.get("$location");
                                $location.path("/login");
                            }
                            else {
                                var $mdToast = $injector.get("$mdToast");
                                $mdToast.show($mdToast.simple()
                                    .textContent(rejection.statusText)
                                    .position("top right"));
                            }
                            return $q.reject(rejection);
                        }
                    };
                }]);
        }]);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaGFwdGljLmNvcmUuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi90cy9jb21wb25lbnRzL2NvcmUvaGFwdGljLmNvcmUudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6IjtJQTJCQSxJQUFJLGFBQWEsR0FBRyxNQUFNLENBQUM7SUFDM0IsSUFBSSxHQUFHLEdBQUcsT0FBTyxDQUFDLE1BQU0sQ0FBQyxTQUFTLEdBQUcsYUFBYSxFQUFFLENBQUMsNEJBQTRCLENBQUMsQ0FBQyxDQUFDO0lBRXBGLEdBQUcsQ0FBQyxNQUFNLENBQUMsQ0FBQyxxQkFBcUIsRUFBRSxVQUFVLEVBQUUsc0JBQXNCLEVBQUUsb0JBQW9CLEVBQUUsNEJBQTRCLEVBQUUsZUFBZSxFQUFFLFVBQzNJLG1CQUFnRCxFQUNoRCxRQUFzQyxFQUN0QyxvQkFBeUIsRUFDekIsa0JBQWlELEVBQ2pELDBCQUF5RCxFQUN6RCxhQUFvQztZQUVwQyxtQ0FBd0IsQ0FBQyxHQUFHLEVBQUUsbUJBQW1CLEVBQUUsUUFBUSxDQUFDLENBQUM7WUFHN0QsMEJBQTBCLENBQUMsVUFBVSxDQUFDLEtBQUssQ0FBQyxDQUFDO1lBRzdDLGtCQUFrQixDQUFDLFNBQVMsQ0FBQyxVQUFTLFNBQXdDLEVBQUUsU0FBbUM7Z0JBQ2xILElBQUksTUFBTSxHQUFHLFFBQVEsQ0FBQztnQkFDdEIsRUFBRSxDQUFDLENBQUMsU0FBUyxDQUFDLEdBQUcsRUFBRSxDQUFDLEtBQUssQ0FBQyxDQUFDLEVBQUUsTUFBTSxDQUFDLE1BQU0sQ0FBQyxLQUFLLE1BQU0sQ0FBQyxDQUFDLENBQUM7b0JBQ3hELE1BQU0sQ0FBQyxNQUFNLEdBQUcsV0FBVyxDQUFDO2dCQUM3QixDQUFDO2dCQUFDLElBQUksQ0FBQyxDQUFDO29CQUNQLE1BQU0sQ0FBQyxHQUFHLENBQUM7Z0JBQ1osQ0FBQztZQUNGLENBQUMsQ0FBQyxDQUFDO1lBR0gsSUFBSSxNQUFNLEdBQXdCLENBQUM7b0JBQ2xDLFFBQVEsRUFBRSxJQUFJO29CQUNkLElBQUksRUFBRSxPQUFPO29CQUNiLEdBQUcsRUFBRSxRQUFRO29CQUNiLFVBQVUsRUFBRSxVQUFVO29CQUN0QixZQUFZLEVBQUUsVUFBVTtvQkFDeEIsV0FBVyxFQUFFLHlCQUFjLENBQUMsYUFBYSxFQUFFLFlBQVksQ0FBQztpQkFDeEQsQ0FBQyxDQUFDO1lBQ0gsbUNBQXdCLENBQUMsYUFBYSxFQUFFLG9CQUFvQixFQUFFLE1BQU0sQ0FBQyxDQUFDO1lBR3RFLEVBQUUsQ0FBQyxDQUFDLFlBQVksQ0FBQyxhQUFhLENBQUMsQ0FBQyxDQUFDLENBQUM7Z0JBQ2pDLGFBQWEsQ0FBQyxRQUFRLENBQUMsT0FBTyxDQUFDLE1BQU0sQ0FBQyxlQUFlLENBQUMsR0FBRyxTQUFTLEdBQUcsWUFBWSxDQUFDLGFBQWEsQ0FBQyxDQUFDO1lBQ2xHLENBQUM7WUFHRCxhQUFhLENBQUMsWUFBWSxDQUFDLElBQUksQ0FBQztnQkFDL0IsTUFBTSxDQUFDO29CQUNOLFNBQVMsRUFBRSxVQUFTLE1BQVc7d0JBQzlCLElBQUksR0FBRyxHQUFHLFFBQVEsQ0FBQyxjQUFjLENBQUMsYUFBYSxDQUFDLENBQUM7d0JBQ2pELEVBQUUsQ0FBQyxDQUFDLEdBQUcsQ0FBQyxDQUFDLENBQUM7NEJBQ1QsR0FBRyxDQUFDLEtBQUssQ0FBQyxVQUFVLEdBQUcsU0FBUyxDQUFDO3dCQUNsQyxDQUFDO3dCQUNELE1BQU0sQ0FBQyxNQUFNLENBQUM7b0JBQ2YsQ0FBQztvQkFDRCxVQUFVLEVBQUUsVUFBUyxRQUFhO3dCQUNqQyxJQUFJLEdBQUcsR0FBRyxRQUFRLENBQUMsY0FBYyxDQUFDLGFBQWEsQ0FBQyxDQUFDO3dCQUNqRCxFQUFFLENBQUMsQ0FBQyxHQUFHLENBQUMsQ0FBQyxDQUFDOzRCQUNULEdBQUcsQ0FBQyxLQUFLLENBQUMsVUFBVSxHQUFHLFFBQVEsQ0FBQzt3QkFDakMsQ0FBQzt3QkFDRCxNQUFNLENBQUMsUUFBUSxDQUFDO29CQUNqQixDQUFDO2lCQUNELENBQUM7WUFDSCxDQUFDLENBQUMsQ0FBQztZQUdILGFBQWEsQ0FBQyxZQUFZLENBQUMsSUFBSSxDQUFDLENBQUMsV0FBVyxFQUFFLElBQUksRUFBRSxVQUFTLFNBQXdDLEVBQUUsRUFBcUI7b0JBQzNILE1BQU0sQ0FBQzt3QkFDTixlQUFlLEVBQUUsVUFBUyxTQUErQzs0QkFDeEUsRUFBRSxDQUFDLENBQUMsU0FBUyxDQUFDLE1BQU0sS0FBSyxHQUFHLElBQUksU0FBUyxDQUFDLE1BQU0sS0FBSyxHQUFHLENBQUMsQ0FBQyxDQUFDO2dDQUMxRCxJQUFJLFNBQVMsR0FBNkIsU0FBUyxDQUFDLEdBQUcsQ0FBQyxXQUFXLENBQUMsQ0FBQztnQ0FDckUsU0FBUyxDQUFDLElBQUksQ0FBQyxRQUFRLENBQUMsQ0FBQzs0QkFDMUIsQ0FBQzs0QkFBQyxJQUFJLENBQUMsQ0FBQztnQ0FDUCxJQUFJLFFBQVEsR0FBbUMsU0FBUyxDQUFDLEdBQUcsQ0FBQyxVQUFVLENBQUMsQ0FBQztnQ0FDekUsUUFBUSxDQUFDLElBQUksQ0FDWixRQUFRLENBQUMsTUFBTSxFQUFFO3FDQUNmLFdBQVcsQ0FBQyxTQUFTLENBQUMsVUFBVSxDQUFDO3FDQUNqQyxRQUFRLENBQUMsV0FBVyxDQUFDLENBQ3ZCLENBQUM7NEJBQ0gsQ0FBQzs0QkFDRCxNQUFNLENBQUMsRUFBRSxDQUFDLE1BQU0sQ0FBQyxTQUFTLENBQUMsQ0FBQzt3QkFDN0IsQ0FBQztxQkFDRCxDQUFDO2dCQUNILENBQUMsQ0FBQyxDQUFDLENBQUM7UUFFTCxDQUFDLENBQUMsQ0FBQyxDQUFDIn0=