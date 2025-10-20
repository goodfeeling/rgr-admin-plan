import apiClient from "../apiClient";

import type { Menu, MenuTreeUserGroup, PageList } from "#/entity";

export class MenuService {
	/**
	 * 获取菜单列表
	 * @param groupId 组ID
	 */
	getMenus(groupId: number) {
		return apiClient.get<Menu[]>({
			url: `${MenuService.Client.Menu}?group_id=${groupId}`,
		});
	}

	/**
	 * 更新菜单
	 * @param id 菜单ID
	 * @param menuInfo 菜单信息
	 */
	updateMenu(id: number, menuInfo: Menu) {
		return apiClient.put<Menu>({
			url: `${MenuService.Client.Menu}/${id}`,
			data: menuInfo,
		});
	}

	/**
	 * 创建菜单
	 * @param menuInfo 菜单信息
	 */
	createMenu(menuInfo: Menu) {
		return apiClient.post<Menu>({
			url: `${MenuService.Client.Menu}`,
			data: menuInfo,
		});
	}

	/**
	 * 搜索菜单分页列表
	 * @param searchStr 搜索字符串
	 */
	searchPageList(searchStr: string) {
		return apiClient.get<PageList<Menu>>({
			url: `${MenuService.Client.SearchMenu}?${searchStr}`,
		});
	}

	/**
	 * 删除菜单
	 * @param id 菜单ID
	 */
	deleteMenu(id: number) {
		return apiClient.delete<string>({
			url: `${MenuService.Client.Menu}/${id}`,
		});
	}

	/**
	 * 获取用户菜单
	 * @param isAll 是否获取所有菜单
	 */
	getUserMenu(isAll = false) {
		return apiClient.get<MenuTreeUserGroup[]>({
			url: `${MenuService.Client.UserMenu}?all=${isAll}`,
		});
	}
}

export namespace MenuService {
	export enum Client {
		Menu = "/menu",
		SearchMenu = "/menu/search",
		UserMenu = "/menu/user",
	}
}

export default new MenuService();
