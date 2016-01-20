define(["require", "exports", "../../core/services/AuthenticationSvc"], function (require, exports) {
    "use strict";
    var LoginCtrl = (function () {
        function LoginCtrl($location, $mdToast, $http, authSvc, $stateParams) {
            this.$location = $location;
            this.$mdToast = $mdToast;
            this.$http = $http;
            this.authSvc = authSvc;
            this.credentials = {
                "email": "",
                "password": ""
            };
            if ($stateParams["logout"]) {
                authSvc.logout();
            }
        }
        LoginCtrl.prototype.signIn = function () {
            var _this = this;
            sessionStorage.setItem("user", this.credentials.email);
            this.authSvc
                .login(this.credentials)
                .then(function () {
                _this.$http.get("/api/me")
                    .success(function (res) {
                    if (res.IsAdmin === true) {
                        _this.$location.path("/admin");
                    }
                    else {
                        _this.$location.path("/");
                    }
                });
            }, function (error) {
                _this.$mdToast.show(_this.$mdToast.simple()
                    .textContent("Authentication failed: Email or Password incorrect")
                    .position("top right"));
            });
        };
        LoginCtrl.$inject = [
            "$location",
            "$mdToast",
            "$http",
            "AuthenticationSvc",
            "$stateParams"
        ];
        return LoginCtrl;
    })();
    exports.LoginCtrl = LoginCtrl;
    angular.module("haptic.login").controller("LoginCtrl", LoginCtrl);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiTG9naW5DdHJsLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vdHMvY29tcG9uZW50cy9sb2dpbi9jb250cm9sbGVycy9Mb2dpbkN0cmwudHMiXSwibmFtZXMiOlsiTG9naW5DdHJsIiwiTG9naW5DdHJsLmNvbnN0cnVjdG9yIiwiTG9naW5DdHJsLnNpZ25JbiJdLCJtYXBwaW5ncyI6IjtJQTBCQSxZQUFZLENBQUM7SUFFYjtRQVlDQSxtQkFDU0EsU0FBbUNBLEVBQ25DQSxRQUF3Q0EsRUFDeENBLEtBQTJCQSxFQUMzQkEsT0FBMEJBLEVBQ2xDQSxZQUE0Q0E7WUFKcENDLGNBQVNBLEdBQVRBLFNBQVNBLENBQTBCQTtZQUNuQ0EsYUFBUUEsR0FBUkEsUUFBUUEsQ0FBZ0NBO1lBQ3hDQSxVQUFLQSxHQUFMQSxLQUFLQSxDQUFzQkE7WUFDM0JBLFlBQU9BLEdBQVBBLE9BQU9BLENBQW1CQTtZQUdsQ0EsSUFBSUEsQ0FBQ0EsV0FBV0EsR0FBR0E7Z0JBQ2xCQSxPQUFPQSxFQUFFQSxFQUFFQTtnQkFDWEEsVUFBVUEsRUFBRUEsRUFBRUE7YUFDZEEsQ0FBQ0E7WUFFRkEsRUFBRUEsQ0FBQ0EsQ0FBQ0EsWUFBWUEsQ0FBQ0EsUUFBUUEsQ0FBQ0EsQ0FBQ0EsQ0FBQ0EsQ0FBQ0E7Z0JBQzVCQSxPQUFPQSxDQUFDQSxNQUFNQSxFQUFFQSxDQUFDQTtZQUNsQkEsQ0FBQ0E7UUFDRkEsQ0FBQ0E7UUFFREQsMEJBQU1BLEdBQU5BO1lBQUFFLGlCQXFCQ0E7WUFwQkFBLGNBQWNBLENBQUNBLE9BQU9BLENBQUNBLE1BQU1BLEVBQUVBLElBQUlBLENBQUNBLFdBQVdBLENBQUNBLEtBQUtBLENBQUNBLENBQUNBO1lBQ3ZEQSxJQUFJQSxDQUFDQSxPQUFPQTtpQkFDVkEsS0FBS0EsQ0FBQ0EsSUFBSUEsQ0FBQ0EsV0FBV0EsQ0FBQ0E7aUJBQ3ZCQSxJQUFJQSxDQUNKQTtnQkFDQ0EsS0FBSUEsQ0FBQ0EsS0FBS0EsQ0FBQ0EsR0FBR0EsQ0FBQ0EsU0FBU0EsQ0FBQ0E7cUJBQ3ZCQSxPQUFPQSxDQUFDQSxVQUFDQSxHQUFRQTtvQkFDakJBLEVBQUVBLENBQUNBLENBQUNBLEdBQUdBLENBQUNBLE9BQU9BLEtBQUtBLElBQUlBLENBQUNBLENBQUNBLENBQUNBO3dCQUMxQkEsS0FBSUEsQ0FBQ0EsU0FBU0EsQ0FBQ0EsSUFBSUEsQ0FBQ0EsUUFBUUEsQ0FBQ0EsQ0FBQ0E7b0JBQy9CQSxDQUFDQTtvQkFBQ0EsSUFBSUEsQ0FBQ0EsQ0FBQ0E7d0JBQ1BBLEtBQUlBLENBQUNBLFNBQVNBLENBQUNBLElBQUlBLENBQUNBLEdBQUdBLENBQUNBLENBQUNBO29CQUMxQkEsQ0FBQ0E7Z0JBQ0ZBLENBQUNBLENBQUNBLENBQUNBO1lBQ0xBLENBQUNBLEVBQ0RBLFVBQUNBLEtBQVVBO2dCQUNWQSxLQUFJQSxDQUFDQSxRQUFRQSxDQUFDQSxJQUFJQSxDQUNqQkEsS0FBSUEsQ0FBQ0EsUUFBUUEsQ0FBQ0EsTUFBTUEsRUFBRUE7cUJBQ3JCQSxXQUFXQSxDQUFDQSxvREFBb0RBLENBQUNBO3FCQUNqRUEsUUFBUUEsQ0FBQ0EsV0FBV0EsQ0FBQ0EsQ0FBQ0EsQ0FBQ0E7WUFDMUJBLENBQUNBLENBQUNBLENBQUNBO1FBQ05BLENBQUNBO1FBOUNNRixpQkFBT0EsR0FBR0E7WUFDaEJBLFdBQVdBO1lBQ1hBLFVBQVVBO1lBQ1ZBLE9BQU9BO1lBQ1BBLG1CQUFtQkE7WUFDbkJBLGNBQWNBO1NBQ2RBLENBQUNBO1FBeUNIQSxnQkFBQ0E7SUFBREEsQ0FBQ0EsQUFuREQsSUFtREM7SUFuRFksaUJBQVMsWUFtRHJCLENBQUE7SUFFRCxPQUFPLENBQUMsTUFBTSxDQUFDLGNBQWMsQ0FBQyxDQUFDLFVBQVUsQ0FBQyxXQUFXLEVBQUUsU0FBUyxDQUFDLENBQUMifQ==