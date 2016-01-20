define(["require", "exports", "../../applications/services/ApplicationsSvc"], function (require, exports) {
    "use strict";
    var PresenterCtrl = (function () {
        function PresenterCtrl($state, appsSvc) {
            this.$state = $state;
            this.appsSvc = appsSvc;
            this.loadApplications();
            this.user = sessionStorage.getItem("user");
        }
        PresenterCtrl.prototype.loadApplications = function () {
            var _this = this;
            return this.appsSvc.getApplicationForUser().then(function (applications) {
                _this.applications = applications;
            });
        };
        PresenterCtrl.prototype.openApplication = function (application, e) {
            var appToken = btoa(application.ConnectionName + "\0c\0noauthlogged");
            var url = "/guacamole/#/client/" + appToken;
            if (localStorage["accessToken"]) {
                url += "?access_token=" + localStorage["accessToken"];
            }
            window.open(url, "_blank");
        };
        PresenterCtrl.prototype.navigateTo = function (loc, e) {
            window.open(loc, "_blank");
        };
        PresenterCtrl.prototype.logout = function () {
            this.$state.go("logout");
        };
        PresenterCtrl.$inject = [
            "$state",
            "ApplicationsSvc"
        ];
        return PresenterCtrl;
    })();
    exports.PresenterCtrl = PresenterCtrl;
    angular.module("haptic.presenter").controller("PresenterCtrl", PresenterCtrl);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiUHJlc2VudGVyQ3RybC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3RzL2NvbXBvbmVudHMvcHJlc2VudGVyL2NvbnRyb2xsZXJzL1ByZXNlbnRlckN0cmwudHMiXSwibmFtZXMiOlsiUHJlc2VudGVyQ3RybCIsIlByZXNlbnRlckN0cmwuY29uc3RydWN0b3IiLCJQcmVzZW50ZXJDdHJsLmxvYWRBcHBsaWNhdGlvbnMiLCJQcmVzZW50ZXJDdHJsLm9wZW5BcHBsaWNhdGlvbiIsIlByZXNlbnRlckN0cmwubmF2aWdhdGVUbyIsIlByZXNlbnRlckN0cmwubG9nb3V0Il0sIm1hcHBpbmdzIjoiO0lBMEJBLFlBQVksQ0FBQztJQUViO1FBVUNBLHVCQUNTQSxNQUFnQ0EsRUFDaENBLE9BQXdCQTtZQUR4QkMsV0FBTUEsR0FBTkEsTUFBTUEsQ0FBMEJBO1lBQ2hDQSxZQUFPQSxHQUFQQSxPQUFPQSxDQUFpQkE7WUFFaENBLElBQUlBLENBQUNBLGdCQUFnQkEsRUFBRUEsQ0FBQ0E7WUFDeEJBLElBQUlBLENBQUNBLElBQUlBLEdBQUdBLGNBQWNBLENBQUNBLE9BQU9BLENBQUNBLE1BQU1BLENBQUNBLENBQUNBO1FBQzVDQSxDQUFDQTtRQUVERCx3Q0FBZ0JBLEdBQWhCQTtZQUFBRSxpQkFJQ0E7WUFIQUEsTUFBTUEsQ0FBQ0EsSUFBSUEsQ0FBQ0EsT0FBT0EsQ0FBQ0EscUJBQXFCQSxFQUFFQSxDQUFDQSxJQUFJQSxDQUFDQSxVQUFDQSxZQUE0QkE7Z0JBQzdFQSxLQUFJQSxDQUFDQSxZQUFZQSxHQUFHQSxZQUFZQSxDQUFDQTtZQUNsQ0EsQ0FBQ0EsQ0FBQ0EsQ0FBQ0E7UUFDSkEsQ0FBQ0E7UUFFREYsdUNBQWVBLEdBQWZBLFVBQWdCQSxXQUF5QkEsRUFBRUEsQ0FBYUE7WUFDdkRHLElBQUlBLFFBQVFBLEdBQUdBLElBQUlBLENBQUNBLFdBQVdBLENBQUNBLGNBQWNBLEdBQUdBLG1CQUFtQkEsQ0FBQ0EsQ0FBQ0E7WUFDdEVBLElBQUlBLEdBQUdBLEdBQUdBLHNCQUFzQkEsR0FBR0EsUUFBUUEsQ0FBQ0E7WUFDNUNBLEVBQUVBLENBQUNBLENBQUNBLFlBQVlBLENBQUNBLGFBQWFBLENBQUNBLENBQUNBLENBQUNBLENBQUNBO2dCQUNqQ0EsR0FBR0EsSUFBSUEsZ0JBQWdCQSxHQUFHQSxZQUFZQSxDQUFDQSxhQUFhQSxDQUFDQSxDQUFDQTtZQUN2REEsQ0FBQ0E7WUFDREEsTUFBTUEsQ0FBQ0EsSUFBSUEsQ0FBQ0EsR0FBR0EsRUFBRUEsUUFBUUEsQ0FBQ0EsQ0FBQ0E7UUFDNUJBLENBQUNBO1FBRURILGtDQUFVQSxHQUFWQSxVQUFXQSxHQUFXQSxFQUFFQSxDQUFhQTtZQUNwQ0ksTUFBTUEsQ0FBQ0EsSUFBSUEsQ0FBQ0EsR0FBR0EsRUFBRUEsUUFBUUEsQ0FBQ0EsQ0FBQ0E7UUFDNUJBLENBQUNBO1FBRURKLDhCQUFNQSxHQUFOQTtZQUNDSyxJQUFJQSxDQUFDQSxNQUFNQSxDQUFDQSxFQUFFQSxDQUFDQSxRQUFRQSxDQUFDQSxDQUFDQTtRQUMxQkEsQ0FBQ0E7UUFsQ01MLHFCQUFPQSxHQUFHQTtZQUNoQkEsUUFBUUE7WUFDUkEsaUJBQWlCQTtTQUNqQkEsQ0FBQ0E7UUFpQ0hBLG9CQUFDQTtJQUFEQSxDQUFDQSxBQXpDRCxJQXlDQztJQXpDWSxxQkFBYSxnQkF5Q3pCLENBQUE7SUFFRCxPQUFPLENBQUMsTUFBTSxDQUFDLGtCQUFrQixDQUFDLENBQUMsVUFBVSxDQUFDLGVBQWUsRUFBRSxhQUFhLENBQUMsQ0FBQyJ9