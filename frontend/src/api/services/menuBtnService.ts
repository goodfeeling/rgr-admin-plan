import apiClient from "../apiClient";

import type { MenuBtn } from "#/entity";

export class MenuBtnService {
	/**
	 * 获取菜单按钮列表
	 * @param menuId 菜单ID
	 */
	getMenuBtns(menuId: number) {
		return apiClient.get<MenuBtn[]>({
			url: `${MenuBtnService.Client.MenuBtn}?menu_id=${menuId}`,
		});
	}

	/**
	 * 更新菜单按钮
	 * @param id 按钮ID
	 * @param apiInfo 按钮信息
	 */
	updateMenuBtn(id: number, apiInfo: MenuBtn) {
		return apiClient.put<MenuBtn>({
			url: `${MenuBtnService.Client.MenuBtn}/${id}`,
			data: apiInfo,
		});
	}

	/**
	 * 创建菜单按钮
	 * @param apiInfo 按钮信息
	 */
	createMenuBtn(apiInfo: MenuBtn) {
		return apiClient.post<MenuBtn>({
			url: `${MenuBtnService.Client.MenuBtn}`,
			data: apiInfo,
		});
	}

	/**
	 * 删除菜单按钮
	 * @param id 按钮ID
	 */
	deleteMenuBtn(id: number) {
		return apiClient.delete<string>({
			url: `${MenuBtnService.Client.MenuBtn}/${id}`,
		});
	}
}

export namespace MenuBtnService {
	export enum Client {
		MenuBtn = "/menu_btn",
	}
}

export default new MenuBtnService();
