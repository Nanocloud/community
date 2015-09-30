/*
 * Nanocloud Community, a comprehensive platform to turn any application
 * into a cloud solution.
 *
 * Copyright (C) 2015 Nanocloud Software
 *
 * This file is part of Nanocloud community.
 *
 * Nanocloud community is free software; you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Nanocloud community is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
