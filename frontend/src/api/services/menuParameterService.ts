import apiClient from "../apiClient";

import type { MenuParameter } from "#/entity";

export class MenuParameterService {
	/**
	 * 获取菜单参数列表
	 * @param menuId 菜单ID
	 */
	getMenuParameters(menuId: number) {
		return apiClient.get<MenuParameter[]>({
			url: `${MenuParameterService.Client.MenuParameter}?menu_id=${menuId}`,
		});
	}

	/**
	 * 更新菜单参数
	 * @param id 参数ID
	 * @param apiInfo 参数信息
	 */
	updateMenuParameter(id: number, apiInfo: MenuParameter) {
		return apiClient.put<MenuParameter>({
			url: `${MenuParameterService.Client.MenuParameter}/${id}`,
			data: apiInfo,
		});
	}

	/**
	 * 创建菜单参数
	 * @param apiInfo 参数信息
	 */
	createMenuParameter(apiInfo: MenuParameter) {
		return apiClient.post<MenuParameter>({
			url: `${MenuParameterService.Client.MenuParameter}`,
			data: apiInfo,
		});
	}

	/**
	 * 删除菜单参数
	 * @param id 参数ID
	 */
	deleteMenuParameter(id: number) {
		return apiClient.delete<string>({
			url: `${MenuParameterService.Client.MenuParameter}/${id}`,
		});
	}
}

export namespace MenuParameterService {
	export enum Client {
		MenuParameter = "/menu_parameter",
	}
}

export default new MenuParameterService();
