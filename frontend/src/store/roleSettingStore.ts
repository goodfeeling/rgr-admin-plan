import roleService from "@/api/services/roleService";
import { create } from "zustand";

interface RoleSettingState {
	ApiIds: string[];
	MenuIds: { [key: string]: number[] };
	RoleBtns: { [key: string]: number[] };
	loading: boolean;
	error: string | null;
	actions: {
		fetch: (roleId: number) => Promise<void>;
		updateMenus: (roleId: number, groupId: string, addMenuId: number[]) => Promise<void>;
		updateApis: (roleId: number, apiIds: string[]) => Promise<void>;
		updateRoleBtns: (roleId: number, menuId: number, btnIds: number[]) => Promise<void>;
		updateRouterPath: (roleId: number, routerPath: string) => Promise<void>;
	};
}
export const useRoleSettingActions = () => useRoleSettingBtnStore((state) => state.actions);
export const useRoleSettingMenuIds = () => useRoleSettingBtnStore((state) => state.MenuIds);
export const useRoleSettingApiIds = () => useRoleSettingBtnStore((state) => state.ApiIds);
export const useRoleSettingBtnIds = () => useRoleSettingBtnStore((state) => state.RoleBtns);
export const useRoleSettingLoading = () => useRoleSettingBtnStore((state) => state.loading);
export const useRoleSettingError = () => useRoleSettingBtnStore((state) => state.error);

const useRoleSettingBtnStore = create<RoleSettingState>()((set, get) => ({
	ApiIds: [],
	MenuIds: {},
	RoleBtns: {},

	loading: true,
	error: null,

	actions: {
		fetch: async (roleId: number) => {
			set({ loading: true, error: null });
			try {
				const response = await roleService.getRoleSetting(roleId);

				set({
					ApiIds: response.role_apis,
					MenuIds: response.role_menus,
					RoleBtns: response.role_btns,
					loading: false,
				});
			} catch (err) {
				console.error(err);
				set({ error: (err as Error).message, loading: false });
			}
		},
		updateMenus: async (roleId: number, groupId: string, addMenuId: number[]) => {
			set({ loading: true, error: null });
			try {
				const menuIds: number[] = [];
				const state = get();
				const menuGroupIds = { ...state.MenuIds };
				// 不存在的数据初始化为空
				if (menuGroupIds[groupId] === undefined) {
					menuGroupIds[groupId] = [];
				}
				for (const gId in menuGroupIds) {
					if (gId === groupId) {
						menuGroupIds[gId] = addMenuId;
					}
					for (const index in menuGroupIds[gId]) {
						menuIds.push(Number(menuGroupIds[gId][index]));
					}
				}
				await roleService.updateRoleMenus(roleId, menuIds);

				set({ MenuIds: menuGroupIds, loading: false });
			} catch (err) {
				console.error(err);

				set({ error: (err as Error).message });
			}
		},
		updateApis: async (roleId: number, apiIds: string[]) => {
			set({ loading: true, error: null });
			try {
				await roleService.updateRoleApis(roleId, apiIds);
				set({ ApiIds: apiIds, loading: false });
			} catch (err) {
				console.error(err);
				set({ error: (err as Error).message });
			}
		},
		updateRouterPath: async (roleId: number, routerPath: string) => {
			set({ loading: true, error: null });
			try {
				await roleService.updateDefaultRouter(roleId, routerPath);
			} catch (err) {
				console.error(err);
			}
		},
		updateRoleBtns: async (roleId: number, menuId: number, btnIds: number[]) => {
			set({ loading: true, error: null });
			try {
				await roleService.updateRoleBtns(roleId, menuId, btnIds);
			} catch (err) {
				console.error(err);
				set({ error: (err as Error).message });
			}
		},
	},
}));
