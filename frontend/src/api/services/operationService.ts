import apiClient from "../apiClient";

import type { Operation, PageList } from "#/entity";

export class OperationService {
	/**
	 * 搜索操作分页列表
	 * @param searchStr 搜索字符串
	 */
	searchPageList(searchStr: string) {
		return apiClient.get<PageList<Operation>>({
			url: `${OperationService.Client.OperationSearch}?${searchStr}`,
		});
	}

	/**
	 * 删除操作
	 * @param id 操作ID
	 */
	deleteOperation(id: number) {
		return apiClient.delete<string>({
			url: `${OperationService.Client.Operation}/${id}`,
		});
	}

	/**
	 * 批量删除操作
	 * @param ids 操作ID数组
	 */
	deleteBatch(ids: number[]) {
		return apiClient.post<number[]>({
			url: OperationService.Client.OperationDeleteBatch,
			data: {
				ids,
			},
		});
	}
}

export namespace OperationService {
	export enum Client {
		Operation = "/operation",
		OperationDeleteBatch = "/operation/delete-batch",
		OperationSearch = "/operation/search",
	}
}

export default new OperationService();
