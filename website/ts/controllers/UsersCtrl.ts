/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
	"use strict";
	
	class UsersCtrl {

		gridOptions: any;

		static $inject = [
			"UserService",
			"$mdDialog"
		];
		constructor(
			private userSrv: UserService,
			private $mdDialog: angular.material.IDialogService
		) {
			this.gridOptions = {
				data: [],
				rowHeight: 36,
				columnDefs: [
					{ field: "Firstname" },
					{ field: "Lastname" },
					{ field: "Email" },
					{
						name: "modifications",
						displayName: "",
						enableColumnMenu: false,
						cellTemplate: "\
							<md-button ng-click='grid.appScope.usersCtrl.startEditUser($event, row.entity)'>\
								<ng-md-icon icon='edit' size='14'></ng-md-icon> Change password\
							</md-button>"
					},
					{
						name: "supression",
						displayName: "",
						enableColumnMenu: false,
						cellTemplate: "\
							<md-button ng-click='grid.appScope.usersCtrl.startDeleteUser($event, row.entity)'>\
								<ng-md-icon icon='delete' size='14'></ng-md-icon> Delete\
							</md-button>"
					}
				]	
			};
			
			this.loadUsers();
		}
		
		get users(): IUser[] {
			return this.gridOptions.data;
		}
		set users(value: IUser[]) {
			this.gridOptions.data = value;
		}
		
		loadUsers(): angular.IPromise<void> {
			return this.userSrv.getAll().then((users: IUser[]) => {
				this.users = users;
			});
		}
		
		startAddUser(e: MouseEvent): angular.IPromise<any> {
			let o = this.getDefaultUserDlgOpt(e);
			o.locals = { user: null };
			return this.$mdDialog
				.show(o);
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
			let i = _.findIndex(this.users, (x: IUser) => x.Email === user.Email);
			if (i >= 0) {
				this.users[i] = user;
			}
		}
		
		startDeleteUser(e: MouseEvent, user: IUser) {
			let o = this.$mdDialog.confirm()
				.parent(angular.element(document.body))
				.title("Delete user")
				.content("Are you sure you want to delete this user?")
				.ok("Yes")
				.cancel("No")
				.targetEvent(e);
			this.$mdDialog
				.show(o)
				.then(this.deleteUser.bind(this, user));
		}
		
		deleteUser(user: IUser) {
			this.userSrv.delete(user);

			// s TODO Haptic does not give user ID for now. We can rely on mail adress for now
			let i = _.findIndex(this.users, (x: IUser) => x.Email === user.Email);
			if (i >= 0) {
				this.users.splice(i, 1);
			}
		}
		
		private getDefaultUserDlgOpt(e: MouseEvent): angular.material.IDialogOptions {
			return {
				controller: "UserCtrl",
				controllerAs: "userCtrl",
				templateUrl: "./views/user.html",
				parent: angular.element(document.body),
				targetEvent: e
			};
		}
	}

	app.controller("UsersCtrl", UsersCtrl);
}
