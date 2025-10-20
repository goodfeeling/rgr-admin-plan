import apiClient from "../apiClient";

import type { MenuGroup, PageList } from "#/entity";

export class MenuGroupService {
	/**
	 * 获取菜单组列表
	 */
	getDictionaries() {
		return apiClient.get<MenuGroup[]>({
			url: `${MenuGroupService.Client.MenuGroup}`,
		});
	}

	/**
	 * 更新菜单组
	 * @param id 菜单组ID
	 * @param apiInfo 菜单组信息
	 */
	updateMenuGroup(id: number, apiInfo: MenuGroup) {
		return apiClient.put<MenuGroup>({
			url: `${MenuGroupService.Client.MenuGroup}/${id}`,
			data: apiInfo,
		});
	}

	/**
	 * 创建菜单组
	 * @param apiInfo 菜单组信息
	 */
	createMenuGroup(apiInfo: MenuGroup) {
		return apiClient.post<MenuGroup>({
			url: `${MenuGroupService.Client.MenuGroup}`,
			data: apiInfo,
		});
	}

	/**
	 * 搜索菜单组分页列表
	 * @param searchStr 搜索字符串
	 */
	searchPageList(searchStr: string) {
		return apiClient.get<PageList<MenuGroup>>({
			url: `${MenuGroupService.Client.SearchMenuGroup}?${searchStr}`,
		});
	}

	/**
	 * 删除菜单组
	 * @param id 菜单组ID
	 */
	deleteMenuGroup(id: number) {
		return apiClient.delete<string>({
			url: `${MenuGroupService.Client.MenuGroup}/${id}`,
		});
	}
}

export namespace MenuGroupService {
	export enum Client {
		MenuGroup = "/menu_group",
		SearchMenuGroup = "/menu_group/search",
	}
}

export default new MenuGroupService();
