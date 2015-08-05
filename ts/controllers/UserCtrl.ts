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
			// here, call the server to save and then should set the id
			if (this.user.Id == null || this.user.Id < 1) { this.user.Id = Math.floor((Math.random() * 999999) + 1); } // fake id
			this.$mdDialog.hide(this.user);
		}
		
		cancel(): void {
			this.$mdDialog.cancel();
		}
		
	}
	
	app.controller("UserCtrl", UserCtrl);	
}
