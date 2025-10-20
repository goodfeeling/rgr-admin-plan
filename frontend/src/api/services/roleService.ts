import apiClient from "../apiClient";

import type { PageList, Role, RoleTree, roleSetting } from "#/entity";

export interface UpdateRole {
	parent_id: number;
	name: string;
	label: string;
	status: number;
	order: number;
	description: string;
}

export class RoleService {
	/**
	 * 获取角色列表
	 * @param status 角色状态
	 */
	getRoles(status = 0) {
		return apiClient.get<Role[]>({
			url: `${RoleService.Client.Role}?status=${status}`,
		});
	}

	/**
	 * 创建角色
	 * @param info 角色信息
	 */
	updateRole(id: number, info: UpdateRole) {
		return apiClient.put<Role>({
			url: `${RoleService.Client.Role}/${id}`,
			data: info,
		});
	}

	/**
	 * 创建角色
	 * @param info 角色信息
	 */
	createRole(info: Role) {
		return apiClient.post<Role>({
			url: `${RoleService.Client.Role}`,
			data: info,
		});
	}

	/**
	 * 获取角色列表
	 * @param searchStr 查询条件
	 */
	searchPageList(searchStr: string) {
		return apiClient.get<PageList<Role>>({
			url: `${RoleService.Client.SearchRole}?${searchStr}`,
		});
	}

	/**
	 * 删除角色
	 * @param id 角色id
	 */
	deleteRole(id: number) {
		return apiClient.delete<string>({
			url: `${RoleService.Client.Role}/${id}`,
		});
	}

	/**
	 * 获取角色树
	 */
	getRoleTree() {
		return apiClient.get<RoleTree>({ url: `${RoleService.Client.RoleTree}` });
	}

	/**
	 * 获取角色设置
	 * @param id 角色id
	 */
	getRoleSetting(id: number) {
		return apiClient.get<roleSetting>({
			url: `${RoleService.Client.Role}/${id}/setting`,
		});
	}

	/**
	 * 更新角色菜单
	 * @param id 角色id
	 * @param menuIds 菜单id
	 */
	updateRoleMenus(id: number, menuIds: number[]) {
		return apiClient.post<boolean>({
			url: `${RoleService.Client.Role}/${id}/menu`,
			data: { menuIds },
		});
	}

	/**
	 * 更新角色接口权限
	 * @param id 角色id
	 * @param apiPaths 接口权限
	 */
	updateRoleApis(id: number, apiPaths: string[]) {
		return apiClient.post<boolean>({
			url: `${RoleService.Client.Role}/${id}/api`,
			data: { apiPaths },
		});
	}

	/**
	 * 绑定按钮权限
	 * @param id 角色id
	 * @param menuId 菜单id
	 * @param btnIds 按钮id列表
	 */
	updateRoleBtns(id: number, menuId: number, btnIds: number[]) {
		return apiClient.post<boolean>({
			url: `${RoleService.Client.Role}/${id}/menu-btns`,
			data: { btnIds, menuId },
		});
	}

	/**
	 * 设置默认路由
	 * @param id 角色id
	 * @param routerPath 路由路径
	 */
	updateDefaultRouter(id: number, routerPath: string) {
		return apiClient.put<Role>({
			url: `${RoleService.Client.Role}/${id}`,
			data: { default_router: routerPath },
		});
	}
}

export namespace RoleService {
	export enum Client {
		Role = "/role",
		SearchRole = "/role/search",
		RoleTree = "/role/tree",
	}
}

export default new RoleService();
