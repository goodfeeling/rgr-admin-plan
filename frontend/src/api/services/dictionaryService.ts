import apiClient from "../apiClient";

import type { Dictionary, PageList } from "#/entity";

export class DictionaryService {
	/**
	 * 获取字典列表
	 */
	getDictionaries() {
		return apiClient.get<Dictionary[]>({
			url: `${DictionaryService.Client.Dictionary}`,
		});
	}

	/**
	 * 更新字典
	 * @param id 字典ID
	 * @param apiInfo 字典信息
	 */
	updateDictionary(id: number, apiInfo: Dictionary) {
		return apiClient.put<Dictionary>({
			url: `${DictionaryService.Client.Dictionary}/${id}`,
			data: apiInfo,
		});
	}

	/**
	 * 创建字典
	 * @param apiInfo 字典信息
	 */
	createDictionary(apiInfo: Dictionary) {
		return apiClient.post<Dictionary>({
			url: `${DictionaryService.Client.Dictionary}`,
			data: apiInfo,
		});
	}

	/**
	 * 搜索字典分页列表
	 * @param searchStr 搜索字符串
	 */
	searchPageList(searchStr: string) {
		return apiClient.get<PageList<Dictionary>>({
			url: `${DictionaryService.Client.SearchDictionary}?${searchStr}`,
		});
	}

	/**
	 * 删除字典
	 * @param id 字典ID
	 */
	deleteDictionary(id: number) {
		return apiClient.delete<string>({
			url: `${DictionaryService.Client.Dictionary}/${id}`,
		});
	}

	/**
	 * 批量删除字典
	 * @param ids 字典ID数组
	 */
	deleteBatch(ids: number[]) {
		return apiClient.post<number>({
			url: `${DictionaryService.Client.DeleteBatch}`,
			data: { ids },
		});
	}

	/**
	 * 根据类型获取字典
	 * @param type 字典类型
	 */
	getByType(type: string) {
		return apiClient.get<Dictionary>({
			url: `${DictionaryService.Client.Type}/${type}`,
		});
	}
}

export namespace DictionaryService {
	export enum Client {
		Dictionary = "/dictionary",
		SearchDictionary = "/dictionary/search",
		GroupsDictionary = "/dictionary/groups",
		DeleteBatch = "/dictionary/batch",
		Type = "/dictionary/type",
	}
}

export default new DictionaryService();
