define(["require", "exports"], function (require, exports) {
    function overrideModuleRegisterer(app, $controllerProvider, $provide) {
        app.controller = function (name, controllerConstructor) {
            $controllerProvider.register(name, controllerConstructor);
            return app;
        };
        app.service = function (name, serviceConstructor) {
            $provide.service(name, serviceConstructor);
            return app;
        };
    }
    exports.overrideModuleRegisterer = overrideModuleRegisterer;
    var requireCtrlStateFactory = ["$q", "futureState", function ($q, futureState) {
            var defer = $q.defer();
            var path = "components/" + futureState.comptName + "/controllers/" + futureState.controller;
            requirejs([path], function () {
                defer.resolve(futureState);
            });
            return defer.promise;
        }];
    function registerCtrlFutureStates(comptName, $futureStateProvider, states) {
        $futureStateProvider.stateFactory("requireCtrl", requireCtrlStateFactory);
        for (var _i = 0; _i < states.length; _i++) {
            var state = states[_i];
            state.type = "requireCtrl";
            state.comptName = comptName;
            $futureStateProvider.futureState(state);
        }
    }
    exports.registerCtrlFutureStates = registerCtrlFutureStates;
    function getTemplateUrl(comptName, url) {
        return "./js/components/" + comptName + "/views/" + url;
    }
    exports.getTemplateUrl = getTemplateUrl;
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiQW1kVG9vbHMuanMiLCJzb3VyY2VSb290IjoiIiwic291cmNlcyI6WyIuLi8uLi8uLi8uLi90cy9jb21wb25lbnRzL2NvcmUvc2VydmljZXMvQW1kVG9vbHMudHMiXSwibmFtZXMiOlsib3ZlcnJpZGVNb2R1bGVSZWdpc3RlcmVyIiwicmVnaXN0ZXJDdHJsRnV0dXJlU3RhdGVzIiwiZ2V0VGVtcGxhdGVVcmwiXSwibWFwcGluZ3MiOiI7SUF5QkEsa0NBQ0MsR0FBb0IsRUFDcEIsbUJBQWdELEVBQ2hELFFBQXNDO1FBRWhDQSxHQUFJQSxDQUFDQSxVQUFVQSxHQUFHQSxVQUFTQSxJQUFZQSxFQUFFQSxxQkFBK0JBO1lBQzdFLG1CQUFtQixDQUFDLFFBQVEsQ0FBQyxJQUFJLEVBQUUscUJBQXFCLENBQUMsQ0FBQztZQUMxRCxNQUFNLENBQUMsR0FBRyxDQUFDO1FBQ1osQ0FBQyxDQUFDQTtRQUVJQSxHQUFJQSxDQUFDQSxPQUFPQSxHQUFHQSxVQUFTQSxJQUFZQSxFQUFFQSxrQkFBNEJBO1lBQ3ZFLFFBQVEsQ0FBQyxPQUFPLENBQUMsSUFBSSxFQUFFLGtCQUFrQixDQUFDLENBQUM7WUFDM0MsTUFBTSxDQUFDLEdBQUcsQ0FBQztRQUNaLENBQUMsQ0FBQ0E7SUFDSEEsQ0FBQ0E7SUFkZSxnQ0FBd0IsMkJBY3ZDLENBQUE7SUFFRCxJQUFJLHVCQUF1QixHQUFHLENBQUMsSUFBSSxFQUFFLGFBQWEsRUFBRSxVQUFTLEVBQXFCLEVBQUUsV0FBOEI7WUFDakgsSUFBSSxLQUFLLEdBQUcsRUFBRSxDQUFDLEtBQUssRUFBRSxDQUFDO1lBQ3ZCLElBQUksSUFBSSxHQUFHLGFBQWEsR0FBUyxXQUFZLENBQUMsU0FBUyxHQUFHLGVBQWUsR0FBRyxXQUFXLENBQUMsVUFBVSxDQUFDO1lBQ25HLFNBQVMsQ0FBQyxDQUFDLElBQUksQ0FBQyxFQUFFO2dCQUNqQixLQUFLLENBQUMsT0FBTyxDQUFDLFdBQVcsQ0FBQyxDQUFDO1lBQzVCLENBQUMsQ0FBQyxDQUFDO1lBQ0gsTUFBTSxDQUFDLEtBQUssQ0FBQyxPQUFPLENBQUM7UUFDdEIsQ0FBQyxDQUFDLENBQUM7SUFHSCxrQ0FBeUMsU0FBaUIsRUFBRSxvQkFBeUIsRUFBRSxNQUEyQjtRQUdqSEMsb0JBQW9CQSxDQUFDQSxZQUFZQSxDQUFDQSxhQUFhQSxFQUFFQSx1QkFBdUJBLENBQUNBLENBQUNBO1FBRTFFQSxHQUFHQSxDQUFDQSxDQUFjQSxVQUFNQSxFQUFuQkEsa0JBQVNBLEVBQVRBLElBQW1CQSxDQUFDQTtZQUFwQkEsSUFBSUEsS0FBS0EsR0FBSUEsTUFBTUEsSUFBVkE7WUFDUEEsS0FBTUEsQ0FBQ0EsSUFBSUEsR0FBR0EsYUFBYUEsQ0FBQ0E7WUFDNUJBLEtBQU1BLENBQUNBLFNBQVNBLEdBQUdBLFNBQVNBLENBQUNBO1lBQ25DQSxvQkFBb0JBLENBQUNBLFdBQVdBLENBQUNBLEtBQUtBLENBQUNBLENBQUNBO1NBQ3hDQTtJQUNGQSxDQUFDQTtJQVZlLGdDQUF3QiwyQkFVdkMsQ0FBQTtJQUdELHdCQUErQixTQUFpQixFQUFFLEdBQVc7UUFDNURDLE1BQU1BLENBQUNBLGtCQUFrQkEsR0FBR0EsU0FBU0EsR0FBR0EsU0FBU0EsR0FBR0EsR0FBR0EsQ0FBQ0E7SUFDekRBLENBQUNBO0lBRmUsc0JBQWMsaUJBRTdCLENBQUEifQ==