import menuService from "@/api/services/menuService";
import type { MenuTreeUserGroup } from "@/types/entity";
import { StorageEnum } from "@/types/enum";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

export type MenuType = MenuTreeUserGroup[];

// 定义 store 状态和方法
interface MenuState {
	menuData: MenuTreeUserGroup[]; // 菜单数据
	loading: boolean; // 加载状态
	error: string | null; // 错误信息
	actions: {
		fetchMenu: () => Promise<void>; // 获取菜单方法
	};
}

const useMenuStore = create<MenuState>()(
	persist(
		(set) => ({
			menuData: [],
			loading: true,
			error: null,

			actions: {
				fetchMenu: async () => {
					set({ loading: true, error: null });
					try {
						const response = await menuService.getUserMenu();
						set({
							menuData: response,
							loading: false,
						});
					} catch (err) {
						set({ error: (err as Error).message, loading: false });
					}
				},
			},
		}),
		{
			name: StorageEnum.Menu, // name of the item in the storage (must be unique)
			storage: createJSONStorage(() => localStorage), // (optional) by default, 'localStorage' is used
			partialize: (state) => ({ [StorageEnum.Menu]: state.menuData }),
		},
	),
);

export const useMenuActions = () => useMenuStore((state) => state.actions);
export const useMenu = () => useMenuStore((state) => state.menuData);

export const useMenuLoading = () => useMenuStore((state) => state.loading);
export const useMenuError = () => useMenuStore((state) => state.error);
