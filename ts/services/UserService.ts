/// <reference path='../../typings/tsd.d.ts' />

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

					for (let usr of JSON.parse(res.result.UsersJsonArray)) {
						users.push(usr);
					}
					return users;
				});
		}

		save(user: IUser): angular.IPromise<boolean> {
			return this.rpc.call({ method: "ServiceUsers.RegisterUser", params: [user], id: 1 })
				.then((res: IRpcResponse): boolean => {
					return this.isError(res);
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
