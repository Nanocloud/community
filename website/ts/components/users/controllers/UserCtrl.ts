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

/// <reference path="../../../../../typings/tsd.d.ts" />
/// <amd-dependency path="../services/UsersSvc" />
import { UsersSvc, IUser } from "../services/UsersSvc";

"use strict";

export class UserCtrl {

	user: IUser;
	userForm: any;
	isCreation: boolean;
	userFormErrorMessage: string;

	static $inject = [
		"UsersSvc",
		"$mdDialog",
		"user"
	];
	constructor(
		private usersSvc: UsersSvc,
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
		if (this.userForm.$invalid || !this.checkPassword(this.user)) {
			return;
		}

		let prm: angular.IPromise<boolean>;
		if (this.isCreation) {
			prm = this.usersSvc.save(this.user);
		} else {
			prm = this.usersSvc.updatePassword(this.user);
		}
		prm.then((ok: boolean) => {
			if (ok) {
				this.$mdDialog.hide(this.user);
			} else {
				this.$mdDialog.cancel();
			}
		});
	}

	checkPassword(user: IUser): boolean {
		// Check if the password meets the following requirements:
		//   * At least 7 characters long
		//   * Less than 65 characters long
		//   * Has at least one digit
		//   * Has at least one Upper case Alphabet
		//   * Has at least one Lower case Alphabet
		// Allowed Characters set:
		//   * Any alphanumeric character 0 to 9 OR A to Z or a to z
		//   * Punctuation symbols [. , " ' ? ! ; : # $ % & ( ) * + - / < > = @ [ ] \ ^ _ { } |]

		if (user === undefined || user === null || typeof user.Password !== "string") {
			return;
		}

		if (user.Password.length < 7 || user.Password.length >= 65) {
			this.userFormErrorMessage = "Password have to contain at least 7 characters and less than 65.";
			return false;
		}

		let oneUpperCaseRegexp = new RegExp("[A-Z]");
		if (!oneUpperCaseRegexp.test(user.Password)) {
			this.userFormErrorMessage = "Password must have at least one upper case letter";
			return false;
		}

		let oneLowerCaseRegexp = new RegExp("[a-z]");
		if (!oneLowerCaseRegexp.test(user.Password)) {
			this.userFormErrorMessage = "Password must have at least one lower case letter";
			return false;
		}

		let oneDigitRegexp = new RegExp("[0-9]");
		if (!oneDigitRegexp.test(user.Password)) {
			this.userFormErrorMessage = "Password must have at least one digit";
			return false;
		}

		let ponctuationMarkRegexp = new RegExp("[^a-zA-Z0-9]");
		if (!ponctuationMarkRegexp.test(user.Password)) {
			this.userFormErrorMessage = "Password must have at least one punctuation mark";
			return false;
		}

		if (user.Password !== user.Password2) {
			this.userFormErrorMessage = "Passwords mismatch.";
			return false;
		}

		return true;
	}

	cancel(): void {
		this.$mdDialog.cancel();
	}

}

angular.module("haptic.users").controller("UserCtrl", UserCtrl);
