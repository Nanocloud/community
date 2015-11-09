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

"use strict";

export interface IRpcCall {
	id?: number | string;
	method: string;
	params?: any;
}

export interface IRpcRequest {
	id?: number | string;
	jsonrpc: string;
	method: string;
	params?: any;
}

export interface IRpcError {
	code: number;
	message: string;
	data: any;
}

export interface IRpcResponse {
	id: number | string;
	jsonrpc: string;
	error: IRpcError;
	result: any;
}

export class RpcSvc {

	static apiUrl = "/rpc";
	static rpcVersion = "2.0";

	static $inject = [
		"$http"
	];
	constructor(
		private $http: angular.IHttpService
	) {

	}
	
	call(data: IRpcCall): angular.IPromise<IRpcResponse> {
		let req = this.getReq(data);
		return this.$http
			.post(RpcSvc.apiUrl, req)
			.then(
				function (res: angular.IHttpPromiseCallbackArg<IRpcResponse>): IRpcResponse {
					return res.data;
				},
				this.xhrToRpcError
			);
	}

	batch(data: IRpcCall[]): angular.IPromise<IRpcResponse[]> {
		let reqs: IRpcRequest[] = [];
		for (let item of data) {
			reqs.push(this.getReq(item));
		}
		return this.$http
			.post(RpcSvc.apiUrl, reqs)
			.then(
				function (res: angular.IHttpPromiseCallbackArg<IRpcResponse[]>): IRpcResponse[] {
					return res.data;
				},
				(res: angular.IHttpPromiseCallbackArg<any>): IRpcResponse[] => {
					return [ this.xhrToRpcError(res) ];
				}
			);
	}

	private getReq(data: IRpcCall): IRpcRequest {
		let req: IRpcRequest = {
			jsonrpc: RpcSvc.rpcVersion,
			method: data.method
		};
		if (data.id !== undefined && data.id !== null) {
			req.id = data.id;
		}
		if (data.params !== undefined && data.params !== null) {
			req.params = data.params;
		} else {
			req.params = [{}];
		}
		return req;
	}

	private xhrToRpcError(res: angular.IHttpPromiseCallbackArg<any>): IRpcResponse {
		if (res.status === 401) {
			window.location.href = "#/login";
		}
		return <IRpcResponse>{
			error: {
				code: res.status,
				message: res.statusText,
				data: res.data
			}
		};
	}

}

angular.module("haptic.core").service("RpcSvc", RpcSvc);
