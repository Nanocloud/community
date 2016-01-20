define(["require", "exports", "../services/ApplicationsSvc"], function (require, exports) {
    "use strict";
    var ApplicationsCtrl = (function () {
        function ApplicationsCtrl(applicationsSrv, $mdDialog) {
            this.applicationsSrv = applicationsSrv;
            this.$mdDialog = $mdDialog;
            this.gridOptions = {
                data: [],
                rowHeight: 36,
                columnDefs: [
                    { field: "ConnectionName" },
                    { field: "Port" },
                    { field: "RemoteApp" },
                    {
                        name: "actions",
                        displayName: "",
                        enableColumnMenu: false,
                        cellTemplate: "\n\t\t\t\t\t\t<md-button ng-click='grid.appScope.applicationsCtrl.openApplication($event, row.entity)'>\n\t\t\t\t\t\t\t<ng-md-icon icon='pageview' size='14'></ng-md-icon> Open\n\t\t\t\t\t\t</md-button>\n\t\t\t\t\t\t<md-button ng-click='grid.appScope.applicationsCtrl.startUnpublishApplication($event, row.entity)'>\n\t\t\t\t\t\t\t<ng-md-icon icon='delete' size='14'></ng-md-icon> Unpublish\n\t\t\t\t\t\t</md-button>"
                    }
                ]
            };
            this.loadApplications();
        }
        Object.defineProperty(ApplicationsCtrl.prototype, "applications", {
            get: function () {
                return this.gridOptions.data;
            },
            set: function (value) {
                this.gridOptions.data = value;
            },
            enumerable: true,
            configurable: true
        });
        ApplicationsCtrl.prototype.loadApplications = function () {
            var _this = this;
            return this.applicationsSrv.getAll().then(function (applications) {
                _this.applications = applications;
            });
        };
        ApplicationsCtrl.prototype.startUnpublishApplication = function (e, application) {
            var o = this.$mdDialog.confirm()
                .parent(angular.element(document.body))
                .title("Unpublish application")
                .content("Are you sure you want to unpublish this application?")
                .ok("Yes")
                .cancel("No")
                .targetEvent(e);
            this.$mdDialog
                .show(o)
                .then(this.unpublishApplication.bind(this, application));
        };
        ApplicationsCtrl.prototype.unpublishApplication = function (application) {
            this.applicationsSrv.unpublish(application);
            var i = _.findIndex(this.applications, function (x) { return x.RemoteApp === application.RemoteApp; });
            if (i >= 0) {
                this.applications.splice(i, 1);
            }
        };
        ApplicationsCtrl.prototype.openApplication = function (e, application) {
            var appToken = btoa(application.ConnectionName + "\0c\0noauthlogged");
            var url = "/guacamole/#/client/" + appToken;
            if (localStorage["accessToken"]) {
                url += "?access_token=" + localStorage["accessToken"];
            }
            window.open(url, "_blank");
        };
        ApplicationsCtrl.prototype.percentDone = function (file) {
            return Math.round(file._prevUploadedSize / file.size * 100).toString();
        };
        ApplicationsCtrl.$inject = [
            "ApplicationsSvc",
            "$mdDialog"
        ];
        return ApplicationsCtrl;
    })();
    exports.ApplicationsCtrl = ApplicationsCtrl;
    angular.module("haptic.applications").controller("ApplicationsCtrl", ApplicationsCtrl);
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiQXBwbGljYXRpb25zQ3RybC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbIi4uLy4uLy4uLy4uL3RzL2NvbXBvbmVudHMvYXBwbGljYXRpb25zL2NvbnRyb2xsZXJzL0FwcGxpY2F0aW9uc0N0cmwudHMiXSwibmFtZXMiOlsiQXBwbGljYXRpb25zQ3RybCIsIkFwcGxpY2F0aW9uc0N0cmwuY29uc3RydWN0b3IiLCJBcHBsaWNhdGlvbnNDdHJsLmFwcGxpY2F0aW9ucyIsIkFwcGxpY2F0aW9uc0N0cmwubG9hZEFwcGxpY2F0aW9ucyIsIkFwcGxpY2F0aW9uc0N0cmwuc3RhcnRVbnB1Ymxpc2hBcHBsaWNhdGlvbiIsIkFwcGxpY2F0aW9uc0N0cmwudW5wdWJsaXNoQXBwbGljYXRpb24iLCJBcHBsaWNhdGlvbnNDdHJsLm9wZW5BcHBsaWNhdGlvbiIsIkFwcGxpY2F0aW9uc0N0cmwucGVyY2VudERvbmUiXSwibWFwcGluZ3MiOiI7SUEwQkEsWUFBWSxDQUFDO0lBRWI7UUFTQ0EsMEJBQ1NBLGVBQWdDQSxFQUNoQ0EsU0FBMENBO1lBRDFDQyxvQkFBZUEsR0FBZkEsZUFBZUEsQ0FBaUJBO1lBQ2hDQSxjQUFTQSxHQUFUQSxTQUFTQSxDQUFpQ0E7WUFFbERBLElBQUlBLENBQUNBLFdBQVdBLEdBQUdBO2dCQUNsQkEsSUFBSUEsRUFBRUEsRUFBRUE7Z0JBQ1JBLFNBQVNBLEVBQUVBLEVBQUVBO2dCQUNiQSxVQUFVQSxFQUFFQTtvQkFDWEEsRUFBRUEsS0FBS0EsRUFBRUEsZ0JBQWdCQSxFQUFFQTtvQkFDM0JBLEVBQUVBLEtBQUtBLEVBQUVBLE1BQU1BLEVBQUVBO29CQUNqQkEsRUFBRUEsS0FBS0EsRUFBRUEsV0FBV0EsRUFBRUE7b0JBQ3RCQTt3QkFDQ0EsSUFBSUEsRUFBRUEsU0FBU0E7d0JBQ2ZBLFdBQVdBLEVBQUVBLEVBQUVBO3dCQUNmQSxnQkFBZ0JBLEVBQUVBLEtBQUtBO3dCQUN2QkEsWUFBWUEsRUFBRUEsaWFBTUFBO3FCQUNkQTtpQkFDREE7YUFD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