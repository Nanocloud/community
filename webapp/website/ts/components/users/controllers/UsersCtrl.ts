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
/// <amd-dependency path="./UserCtrl" />
/// <amd-dependency path="../../services/services/ServicesFct" />
import { UsersSvc, IUser } from "../services/UsersSvc";
import { ServicesFct } from "../../services/services/ServicesFct";

"use strict";

export class UsersCtrl {

	users: any;
	displayHelp: boolean;
	windowsState: boolean = false;
	loadWindowHasFinished: boolean = false;

	static $inject = [
		"UsersSvc",
		"ServicesFct",
		"$mdDialog"
	];
	constructor(
		private usersSvc: UsersSvc,
		private servicesFct: ServicesFct,
		private $mdDialog: angular.material.IDialogService
	) {
		this.loadWindowHasFinished = false;
		this.servicesFct.getWindowsStatus().then((windowsState: boolean) => {
			this.loadWindowHasFinished = true;
			this.windowsState = windowsState;
		});
		this.loadUsers();
		this.displayHelp = false;
	}

	loadUsers(): angular.IPromise<void> {
		return this.usersSvc.getAll().then((users: IUser[]) => {
			this.users = users;
		});
	}

	startAddUser(e: MouseEvent): angular.IPromise<any> {
		let o = this.getDefaultUserDlgOpt(e);
		o.locals = { user: null };
		return this.$mdDialog
			.show(o)
			.then((user: IUser) => {
				if (user) {
					this.users.push(user);
				}
			});
	}
	
	addUser(user: IUser): void {
		this.users.push(user);
	}
	
	startEditUser(e: MouseEvent, user: IUser) {
		let o = this.getDefaultUserDlgOpt(e);
		o.locals = { user: user };
		return this.$mdDialog
			.show(o)
			.then(this.editUser.bind(this));
	}
	
	editUser(user: IUser) {
		// here, call the server to edit
		let i = _.findIndex(this.users, (x: IUser) => x.email === user.email);
		if (i >= 0) {
			this.users[i] = user;
		}
	}
	
	startDeleteUser(e: MouseEvent, user: IUser) {
		let o = this.$mdDialog.confirm()
			.parent(angular.element(document.body))
			.title("Delete user")
			.textContent("Are you sure you want to delete \"" + user.email + "\"?")
			.ok("Yes")
			.cancel("No")
			.targetEvent(e);
		this.$mdDialog
			.show(o)
			.then(this.deleteUser.bind(this, user));
	}
	
	deleteUser(user: IUser) {
		this.usersSvc.delete(user);

		// TODO Haptic does not give user ID for now. We can rely on mail adress for now
		let i = _.findIndex(this.users, (x: IUser) => x.email === user.email);
		if (i >= 0) {
			this.users.splice(i, 1);
		}
	}

	toggleHelp(e: MouseEvent) {
		if (this.displayHelp === true) {
			this.displayHelp = false;
		} else {
			this.displayHelp = true;
		}
	}

	private getDefaultUserDlgOpt(e: MouseEvent): angular.material.IDialogOptions {
		return {
			controller: "UserCtrl",
			controllerAs: "userCtrl",
			templateUrl: "./js/components/users/views/user.html",
			parent: angular.element(document.body),
			targetEvent: e
		};
	}
}

angular.module("haptic.users").controller("UsersCtrl", UsersCtrl);
