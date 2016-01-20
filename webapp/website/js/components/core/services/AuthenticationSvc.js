define(["require", "exports"], function (require, exports) {
    "use strict";
    var AuthenticationSvc = (function () {
        function AuthenticationSvc($http, $q) {
            this.$http = $http;
            this.$q = $q;
        }
        AuthenticationSvc.prototype.login = function (credentials) {
            var _this = this;
            var appKey = "9405fb6b0e59d2997e3c777a22d8f0e617a9f5b36b6565c7579e5be6deb8f7ae";
            var appSecret = "9050d67c2be0943f2c63507052ddedb3ae34a30e39bbbbdab241c93f8b5cf341";
            var basic = btoa(appKey + ":" + appSecret);
            return this.$http
                .post("/oauth/token", JSON.stringify({
                "username": credentials.email,
                "password": credentials.password,
                "grant_type": "password"
            }), {
                headers: { "Authorization": "Basic " + basic }
            })
                .then(function (res) {
                localStorage["accessToken"] = res.data.access_token;
                _this.$http.defaults.headers.common["Authorization"] = "Bearer " + res.data.access_token;
            });
        };
        AuthenticationSvc.prototype.logout = function () {
            var dfr = this.$q.defer();
            sessionStorage.clear();
            localStorage.clear();
            this.$http.defaults.headers.common["Authorization"] = "";
            dfr.resolve();
            return dfr.promise;
        };
        AuthenticationSvc.$inject = [
            "$http",
            "$q"
        ];
        return AuthenticationSvc;
    })();
    exports.AuthenticationSvc = AuthenticationSvc;
    angular.module("haptic.core").service("AuthenticationSvc", AuthenticationSvc);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiQXV0aGVudGljYXRpb25TdmMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi90cy9jb21wb25lbnRzL2NvcmUvc2VydmljZXMvQXV0aGVudGljYXRpb25TdmMudHMiXSwibmFtZXMiOlsiQXV0aGVudGljYXRpb25TdmMiLCJBdXRoZW50aWNhdGlvblN2Yy5jb25zdHJ1Y3RvciIsIkF1dGhlbnRpY2F0aW9uU3ZjLmxvZ2luIiwiQXV0aGVudGljYXRpb25TdmMubG9nb3V0Il0sIm1hcHBpbmdzIjoiO0lBeUJBLFlBQVksQ0FBQztJQUViO1FBTUNBLDJCQUNTQSxLQUEyQkEsRUFDM0JBLEVBQXFCQTtZQURyQkMsVUFBS0EsR0FBTEEsS0FBS0EsQ0FBc0JBO1lBQzNCQSxPQUFFQSxHQUFGQSxFQUFFQSxDQUFtQkE7UUFFOUJBLENBQUNBO1FBRURELGlDQUFLQSxHQUFMQSxVQUFNQSxXQUFnQkE7WUFBdEJFLGlCQWlCQ0E7WUFoQkFBLElBQUlBLE1BQU1BLEdBQUdBLGtFQUFrRUEsQ0FBQ0E7WUFDaEZBLElBQUlBLFNBQVNBLEdBQUdBLGtFQUFrRUEsQ0FBQ0E7WUFFbkZBLElBQUlBLEtBQUtBLEdBQUdBLElBQUlBLENBQUNBLE1BQU1BLEdBQUdBLEdBQUdBLEdBQUdBLFNBQVNBLENBQUNBLENBQUNBO1lBQzNDQSxNQUFNQSxDQUFDQSxJQUFJQSxDQUFDQSxLQUFLQTtpQkFDZkEsSUFBSUEsQ0FBQ0EsY0FBY0EsRUFBRUEsSUFBSUEsQ0FBQ0EsU0FBU0EsQ0FBQ0E7Z0JBQ3BDQSxVQUFVQSxFQUFFQSxXQUFXQSxDQUFDQSxLQUFLQTtnQkFDN0JBLFVBQVVBLEVBQUVBLFdBQVdBLENBQUNBLFFBQVFBO2dCQUNoQ0EsWUFBWUEsRUFBRUEsVUFBVUE7YUFDeEJBLENBQUNBLEVBQUVBO2dCQUNIQSxPQUFPQSxFQUFFQSxFQUFFQSxlQUFlQSxFQUFFQSxRQUFRQSxHQUFHQSxLQUFLQSxFQUFFQTthQUM5Q0EsQ0FBQ0E7aUJBQ0RBLElBQUlBLENBQUNBLFVBQUNBLEdBQXlDQTtnQkFDL0NBLFlBQVlBLENBQUNBLGFBQWFBLENBQUNBLEdBQUdBLEdBQUdBLENBQUNBLElBQUlBLENBQUNBLFlBQVlBLENBQUNBO2dCQUNwREEsS0FBSUEsQ0FBQ0EsS0FBS0EsQ0FBQ0EsUUFBUUEsQ0FBQ0EsT0FBT0EsQ0FBQ0EsTUFBTUEsQ0FBQ0EsZUFBZUEsQ0FBQ0EsR0FBR0EsU0FBU0EsR0FBR0EsR0FBR0EsQ0FBQ0EsSUFBSUEsQ0FBQ0EsWUFBWUEsQ0FBQ0E7WUFDekZBLENBQUNBLENBQUNBLENBQUNBO1FBQ0xBLENBQUNBO1FBRURGLGtDQUFNQSxHQUFOQTtZQUNDRyxJQUFJQSxHQUFHQSxHQUFHQSxJQUFJQSxDQUFDQSxFQUFFQSxDQUFDQSxLQUFLQSxFQUFRQSxDQUFDQTtZQUNoQ0EsY0FBY0EsQ0FBQ0EsS0FBS0EsRUFBRUEsQ0FBQ0E7WUFDdkJBLFlBQVlBLENBQUNBLEtBQUtBLEVBQUVBLENBQUNBO1lBQ3JCQSxJQUFJQSxDQUFDQSxLQUFLQSxDQUFDQSxRQUFRQSxDQUFDQSxPQUFPQSxDQUFDQSxNQUFNQSxDQUFDQSxlQUFlQSxDQUFDQSxHQUFHQSxFQUFFQSxDQUFDQTtZQUN6REEsR0FBR0EsQ0FBQ0EsT0FBT0EsRUFBRUEsQ0FBQ0E7WUFDZEEsTUFBTUEsQ0FBQ0EsR0FBR0EsQ0FBQ0EsT0FBT0EsQ0FBQ0E7UUFDcEJBLENBQUNBO1FBcENNSCx5QkFBT0EsR0FBR0E7WUFDaEJBLE9BQU9BO1lBQ1BBLElBQUlBO1NBQ0pBLENBQUNBO1FBbUNIQSx3QkFBQ0E7SUFBREEsQ0FBQ0EsQUF4Q0QsSUF3Q0M7SUF4Q1kseUJBQWlCLG9CQXdDN0IsQ0FBQTtJQUVELE9BQU8sQ0FBQyxNQUFNLENBQUMsYUFBYSxDQUFDLENBQUMsT0FBTyxDQUFDLG1CQUFtQixFQUFFLGlCQUFpQixDQUFDLENBQUMifQ==