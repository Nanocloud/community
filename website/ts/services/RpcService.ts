/// <reference path='../../../typings/tsd.d.ts' />

module hapticFrontend {
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
	
	export class RpcService {

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
				.post(RpcService.apiUrl, req)
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
				.post(RpcService.apiUrl, reqs)
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
				jsonrpc: RpcService.rpcVersion,
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
				window.location.href = "/login.html";
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

	app.service("RpcService", RpcService);
}
