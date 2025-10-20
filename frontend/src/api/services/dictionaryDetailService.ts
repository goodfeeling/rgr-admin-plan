import apiClient from "../apiClient";

import type { DictionaryDetail, PageList } from "#/entity";

export class DictionaryDetailService {
	/**
	 * 获取字典详情列表
	 */
	getDictionaryDetails() {
		return apiClient.get<DictionaryDetail[]>({
			url: `${DictionaryDetailService.Client.DictionaryDetail}`,
		});
	}

	/**
	 * 更新字典详情
	 * @param id 字典详情ID
	 * @param apiInfo 字典详情信息
	 */
	updateDictionaryDetail(id: number, apiInfo: DictionaryDetail) {
		return apiClient.put<DictionaryDetail>({
			url: `${DictionaryDetailService.Client.DictionaryDetail}/${id}`,
			data: apiInfo,
		});
	}

	/**
	 * 创建字典详情
	 * @param apiInfo 字典详情信息
	 */
	createDictionaryDetail(apiInfo: DictionaryDetail) {
		return apiClient.post<DictionaryDetail>({
			url: `${DictionaryDetailService.Client.DictionaryDetail}`,
			data: apiInfo,
		});
	}

	/**
	 * 搜索字典详情分页列表
	 * @param searchStr 搜索字符串
	 */
	searchPageList(searchStr: string) {
		return apiClient.get<PageList<DictionaryDetail>>({
			url: `${DictionaryDetailService.Client.SearchDictionaryDetail}?${searchStr}`,
		});
	}

	/**
	 * 删除字典详情
	 * @param id 字典详情ID
	 */
	deleteDictionaryDetail(id: number) {
		return apiClient.delete<string>({
			url: `${DictionaryDetailService.Client.DictionaryDetail}/${id}`,
		});
	}

	/**
	 * 批量删除字典详情
	 * @param ids 字典详情ID数组
	 */
	deleteBatch(ids: number[]) {
		return apiClient.post<number>({
			url: `${DictionaryDetailService.Client.DeleteBatch}`,
			data: { ids },
		});
	}
}

export namespace DictionaryDetailService {
	export enum Client {
		DictionaryDetail = "/dictionary_detail",
		SearchDictionaryDetail = "/dictionary_detail/search",
		GroupsDictionaryDetail = "/dictionary_detail/groups",
		DeleteBatch = "/dictionary_detail/batch",
	}
}

export default new DictionaryDetailService();
