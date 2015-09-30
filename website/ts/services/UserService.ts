/*
 * Nanocloud community -- transform any application into SaaS solution
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

	export interface IUser {
		Firstname: string;
		Lastname: string;
		Email: string;
		Password?: string;
		Password2?: string;
	}

	export class UserService {

		static $inject = [
			"RpcService",
			"$mdToast"
		];
		constructor(
			private rpc: RpcService,
			private $mdToast: angular.material.IToastService
		) {

		}

		getAll(): angular.IPromise<IUser[]> {
			return this.rpc.call({ method: "ServiceUsers.GetList", id: 1 })
				.then((res: IRpcResponse): IUser[] => {
					let users: IUser[] = [];

					if (this.isError(res)) {
						return [];
					}

					if (res.result.Users) {
						for (let u of res.result.Users) {
							users.push(u);
						}
					}
					return users;
				});
		}

		save(user: IUser): angular.IPromise<boolean> {
			return this.rpc.call({ method: "ServiceUsers.RegisterUser", params: [user], id: 1 })
				.then((res: IRpcResponse): boolean => {
					return !this.isError(res);
				});
		}

		delete(user: IUser): angular.IPromise<void> {
			return this.rpc.call({ method: "ServiceUsers.DeleteUser", params: [{"Email": user.Email}], id: 1 })
				.then((res: IRpcResponse): void => {
					this.isError(res);
				});
		}

		updatePassword(user: IUser): angular.IPromise<boolean> {
			return this.rpc.call({
				method: "ServiceUsers.UpdateUserPassword",
				params: [{
					"Email": user.Email,
					"Password": user.Password
				}],
				id: 1
			}).then((res: IRpcResponse): boolean => {
				return this.isError(res);
			});
		}

		private isError(res: IRpcResponse): boolean {
			if (res.error == null) {
				return false;
			}
			this.$mdToast.show(
				this.$mdToast.simple()
					.content(res.error.code === 0 ? "Internal Error" : JSON.stringify(res.error))
					.position("top right")
			);
			return true;
		}
		
	}

	app.service("UserService", UserService);
}
