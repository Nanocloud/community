/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";
	
	class UserCtrl {
		
		user: IUser;
		
		static $inject = [
			"UserService",
			"$mdDialog",
			"user"
		];
		constructor(
			private userSrv: UserService,
			private $mdDialog: angular.material.IDialogService,
			user: IUser
		) {
			if (user) {
				this.user = angular.copy(user);
			}
		}
		
		save(): void {
			if (this.userSrv.save(this.user)) {
				this.$mdDialog.hide(this.user);
			}
			this.$mdDialog.cancel();
		}
		
		cancel(): void {
			this.$mdDialog.cancel();
		}
		
	}

	app.controller("UserCtrl", UserCtrl);
}
