/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	class UserCtrl {

		user: IUser;
		userForm: any;
		isCreation: boolean;

		static $inject = [
			"UserService",
			"$mdDialog",
			"user"
		];
		constructor(
			private userSrv: UserService,
			private $mdDialog: angular.material.IDialogService,
			user: IUser,
			isCreation: boolean
		) {
			if (user) {
				this.user = angular.copy(user);
				this.isCreation = false;
			} else {
				this.isCreation = true;
			}
		}

		save(): void {
			let success;

			if (this.userForm.$invalid) {
				return;
			} else if (this.user.Password !== this.user.Password2) {
				return;
			}

			if (this.isCreation) {
				success = this.userSrv.save(this.user);
			} else {
				success = this.userSrv.updatePassword(this.user);
			}
			if (success) {
				this.$mdDialog.hide(this.user);
			} else {
				this.$mdDialog.cancel();
			}
		}
		
		cancel(): void {
			this.$mdDialog.cancel();
		}
		
	}

	app.controller("UserCtrl", UserCtrl);
}
