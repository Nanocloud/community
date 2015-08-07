/// <reference path='../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";

	export interface IUser {
		Id: number;
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
			return this.rpc.call({ method: "UserService.GetUsersList", id: 1 })
				.then((res: IRpcResponse): IUser[] => {
					let users: IUser[] = [];
					
					if (this.isError(res)) {
						// return users;
						return [{ // fake data
							Id: 1,
							Firstname: "John",
							Lastname: "Doe",
							Email: "jdoe@nanoloud.com"
						}];
					}
					
					for (let usr of res.result.UsersJsonArray) {
						users.push(JSON.parse(usr));
					}
					return users;
				});
		}

		save(user: IUser): angular.IPromise<void> {
			return this.rpc.call({ method: "UserService.SaveUser", params: user, id: 1 })
				.then((res: IRpcResponse): void => {
					this.isError(res);
				});
		}

		private isError(res: IRpcResponse): boolean {
			if (res.error == null) {
				return false;
			}
			this.$mdToast.show(
				this.$mdToast.simple()
					.content(res.error.code === 0 ? "Internal Error" : res.error.message)
					.position("top right")
			);
			return true;
		}
		
	}

	app.service("UserService", UserService);
}
