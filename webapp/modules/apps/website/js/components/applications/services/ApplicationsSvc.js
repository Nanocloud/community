define(["require", "exports"], function (require, exports) {
    "use strict";
    var ApplicationsSvc = (function () {
        function ApplicationsSvc($http, $mdToast) {
            this.$http = $http;
            this.$mdToast = $mdToast;
        }
        ApplicationsSvc.prototype.getAll = function () {
            var _this = this;
            return this.$http.get("/api/apps")
                .then(function (res) {
                var apps = res.data || [];
                for (var _i = 0; _i < apps.length; _i++) {
                    var app = apps[_i];
                    app.RemoteApp = _this.cleanAppName(app.RemoteApp);
                }
                return apps;
            }, function () { return []; });
        };
        ApplicationsSvc.prototype.getApplicationForUser = function () {
            var _this = this;
            return this.$http.get("/api/apps/me")
                .then(function (res) {
                var apps = res.data || [];
                for (var _i = 0; _i < apps.length; _i++) {
                    var app = apps[_i];
                    app.RemoteApp = _this.cleanAppName(app.RemoteApp);
                }
                return apps;
            }, function () { return []; });
        };
        ApplicationsSvc.prototype.unpublish = function (application) {
            return this.$http.delete("/api/apps/" + application.RemoteApp);
        };
        ApplicationsSvc.prototype.cleanAppName = function (appName) {
            if (appName) {
                return appName.replace(/^\|\|/, "");
            }
            else {
                return "Desktop";
            }
        };
        ApplicationsSvc.$inject = [
            "$http",
            "$mdToast"
        ];
        return ApplicationsSvc;
    })();
    exports.ApplicationsSvc = ApplicationsSvc;
    angular.module("haptic.applications").service("ApplicationsSvc", ApplicationsSvc);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiQXBwbGljYXRpb25zU3ZjLmpzIiwic291cmNlUm9vdCI6IiIsInNvdXJjZXMiOlsiLi4vLi4vLi4vLi4vdHMvY29tcG9uZW50cy9hcHBsaWNhdGlvbnMvc2VydmljZXMvQXBwbGljYXRpb25zU3ZjLnRzIl0sIm5hbWVzIjpbIkFwcGxpY2F0aW9uc1N2YyIsIkFwcGxpY2F0aW9uc1N2Yy5jb25zdHJ1Y3RvciIsIkFwcGxpY2F0aW9uc1N2Yy5nZXRBbGwiLCJBcHBsaWNhdGlvbnNTdmMuZ2V0QXBwbGljYXRpb25Gb3JVc2VyIiwiQXBwbGljYXRpb25zU3ZjLnVucHVibGlzaCIsIkFwcGxpY2F0aW9uc1N2Yy5jbGVhbkFwcE5hbWUiXSwibWFwcGluZ3MiOiI7SUF3QkEsWUFBWSxDQUFDO0lBV2I7UUFNQ0EseUJBQ1NBLEtBQTJCQSxFQUMzQkEsUUFBd0NBO1lBRHhDQyxVQUFLQSxHQUFMQSxLQUFLQSxDQUFzQkE7WUFDM0JBLGFBQVFBLEdBQVJBLFFBQVFBLENBQWdDQTtRQUdqREEsQ0FBQ0E7UUFFREQsZ0NBQU1BLEdBQU5BO1lBQUFFLGlCQVlDQTtZQVhBQSxNQUFNQSxDQUFDQSxJQUFJQSxDQUFDQSxLQUFLQSxDQUFDQSxHQUFHQSxDQUFDQSxXQUFXQSxDQUFDQTtpQkFDaENBLElBQUlBLENBQ0pBLFVBQUNBLEdBQW9EQTtnQkFDcERBLElBQUlBLElBQUlBLEdBQUdBLEdBQUdBLENBQUNBLElBQUlBLElBQUlBLEVBQUVBLENBQUNBO2dCQUMxQkEsR0FBR0EsQ0FBQ0EsQ0FBWUEsVUFBSUEsRUFBZkEsZ0JBQU9BLEVBQVBBLElBQWVBLENBQUNBO29CQUFoQkEsSUFBSUEsR0FBR0EsR0FBSUEsSUFBSUEsSUFBUkE7b0JBQ1hBLEdBQUdBLENBQUNBLFNBQVNBLEdBQUdBLEtBQUlBLENBQUNBLFlBQVlBLENBQUNBLEdBQUdBLENBQUNBLFNBQVNBLENBQUNBLENBQUNBO2lCQUNqREE7Z0JBQ0RBLE1BQU1BLENBQUNBLElBQUlBLENBQUNBO1lBQ2JBLENBQUNBLEVBQ0RBLGNBQU1BLE9BQUFBLEVBQUVBLEVBQUZBLENBQUVBLENBQ1JBLENBQUNBO1FBQ0pBLENBQUNBO1FBRURGLCtDQUFxQkEsR0FBckJBO1lBQUFHLGlCQVlDQTtZQVhBQSxNQUFNQSxDQUFDQSxJQUFJQSxDQUFDQSxLQUFLQSxDQUFDQSxHQUFHQSxDQUFDQSxjQUFjQSxDQUFDQTtpQkFDbkNBLElBQUlBLENBQ0pBLFVBQUNBLEdBQW9EQTtnQkFDcERBLElBQUlBLElBQUlBLEdBQUdBLEdBQUdBLENBQUNBLElBQUlBLElBQUlBLEVBQUVBLENBQUNBO2dCQUMxQkEsR0FBR0EsQ0FBQ0EsQ0FBWUEsVUFBSUEsRUFBZkEsZ0JBQU9BLEVBQVBBLElBQWVBLENBQUNBO29CQUFoQkEsSUFBSUEsR0FBR0EsR0FBSUEsSUFBSUEsSUFBUkE7b0JBQ1hBLEdBQUdBLENBQUNBLFNBQVNBLEdBQUdBLEtBQUlBLENBQUNBLFlBQVlBLENBQUNBLEdBQUdBLENBQUNBLFNBQVNBLENBQUNBLENBQUNBO2lCQUNqREE7Z0JBQ0RBLE1BQU1BLENBQUNBLElBQUlBLENBQUNBO1lBQ2JBLENBQUNBLEVBQ0RBLGNBQU1BLE9BQUFBLEVBQUVBLEVBQUZBLENBQUVBLENBQ1JBLENBQUNBO1FBQ0pBLENBQUNBO1FBRURILG1DQUFTQSxHQUFUQSxVQUFVQSxXQUF5QkE7WUFDbENJLE1BQU1BLENBQUNBLElBQUlBLENBQUNBLEtBQUtBLENBQUNBLE1BQU1BLENBQUNBLFlBQVlBLEdBQUdBLFdBQVdBLENBQUNBLFNBQVNBLENBQUNBLENBQUNBO1FBQ2hFQSxDQUFDQTtRQUVPSixzQ0FBWUEsR0FBcEJBLFVBQXFCQSxPQUFlQTtZQUNuQ0ssRUFBRUEsQ0FBQ0EsQ0FBQ0EsT0FBT0EsQ0FBQ0EsQ0FBQ0EsQ0FBQ0E7Z0JBQ2JBLE1BQU1BLENBQUNBLE9BQU9BLENBQUNBLE9BQU9BLENBQUNBLE9BQU9BLEVBQUVBLEVBQUVBLENBQUNBLENBQUNBO1lBQ3JDQSxDQUFDQTtZQUFDQSxJQUFJQSxDQUFDQSxDQUFDQTtnQkFDUEEsTUFBTUEsQ0FBQ0EsU0FBU0EsQ0FBQ0E7WUFDbEJBLENBQUNBO1FBQ0ZBLENBQUNBO1FBakRNTCx1QkFBT0EsR0FBR0E7WUFDaEJBLE9BQU9BO1lBQ1BBLFVBQVVBO1NBQ1ZBLENBQUNBO1FBK0NIQSxzQkFBQ0E7SUFBREEsQ0FBQ0EsQUFwREQsSUFvREM7SUFwRFksdUJBQWUsa0JBb0QzQixDQUFBO0lBRUQsT0FBTyxDQUFDLE1BQU0sQ0FBQyxxQkFBcUIsQ0FBQyxDQUFDLE9BQU8sQ0FBQyxpQkFBaUIsRUFBRSxlQUFlLENBQUMsQ0FBQyJ9