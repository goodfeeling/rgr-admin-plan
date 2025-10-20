import apiClient from "../apiClient";

import type { Api, ApiGroupItem, PageList } from "#/entity";

export class ApisService {
	/**
	 * 获取所有API列表
	 */
	getApis() {
		return apiClient.get<Api[]>({ url: `${ApisService.Client.Api}` });
	}

	/**
	 * 更新API信息
	 * @param id API ID
	 * @param apiInfo API信息
	 */
	updateApi(id: number, apiInfo: Api) {
		return apiClient.put<Api>({
			url: `${ApisService.Client.Api}/${id}`,
			data: apiInfo,
		});
	}

	/**
	 * 创建新API
	 * @param apiInfo API信息
	 */
	createApi(apiInfo: Api) {
		return apiClient.post<Api>({
			url: `${ApisService.Client.Api}`,
			data: apiInfo,
		});
	}

	/**
	 * 搜索API分页列表
	 * @param searchStr 搜索字符串
	 */
	searchPageList(searchStr: string) {
		return apiClient.get<PageList<Api>>({
			url: `${ApisService.Client.SearchApi}?${searchStr}`,
		});
	}

	/**
	 * 删除指定API
	 * @param id API ID
	 */
	deleteApi(id: number) {
		return apiClient.delete<string>({ url: `${ApisService.Client.Api}/${id}` });
	}

	/**
	 * 批量删除API
	 * @param ids API ID数组
	 */
	deleteBatch(ids: number[]) {
		return apiClient.post<number>({
			url: `${ApisService.Client.DeleteBatch}`,
			data: { ids },
		});
	}

	/**
	 * 获取API分组列表
	 */
	getApiGroupList(path?: string) {
		return apiClient.get<ApiGroupItem[]>({
			url: `${ApisService.Client.ApiGroupList}?path=${path}`,
		});
	}

	/**
	 * 同步API
	 */
	synchronizeApi() {
		return apiClient.post<{ count: number }>({
			url: ApisService.Client.ApiSynchronize,
		});
	}

	/**
	 * 下载导入模板
	 */
	downloadTemplate() {
		return apiClient.download({
			url: `${ApisService.Client.DownloadTemplate}`,
			method: "GET",
		});
	}

	/**
	 * 导入
	 */
	importApi(file: File) {
		const formData = new FormData();
		formData.append("file", file);
		return apiClient.post<{ count: number }>({
			url: ApisService.Client.Import,
			data: formData,
			headers: {
				"Content-Type": "multipart/form-data",
			},
		});
	}

	/**
	 * 导出
	 */
	exportApi() {
		return apiClient.download({
			url: `${ApisService.Client.Export}`,
			method: "GET",
		});
	}
}

export namespace ApisService {
	export enum Client {
		Api = "/api",
		SearchApi = "/api/search",
		GroupsApi = "/api/groups",
		DeleteBatch = "/api/batch",
		ApiGroupList = "/api/group-list",
		ApiSynchronize = "/api/synchronize",
		DownloadTemplate = "/api/excel/template",
		Import = "/api/excel/import",
		Export = "/api/excel/export",
	}
}

export default new ApisService();
