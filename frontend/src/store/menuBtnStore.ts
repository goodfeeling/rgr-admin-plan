import menuBtnService from "@/api/services/menuBtnService";
import type { MenuBtn } from "@/types/entity";
import { create } from "zustand";

interface MenuBtnState {
	menuBtnData: MenuBtn[];
	loading: boolean;
	error: string | null;
	actions: {
		fetchMenu: (menuID: number) => Promise<void>;
		updateOrCreateMenu: (data: MenuBtn) => Promise<void>;
		deleteMenu: (id: number) => Promise<void>;
	};
}

const useMenuBtnStore = create<MenuBtnState>()((set) => ({
	menuBtnData: [],
	loading: true,
	error: null,

	actions: {
		fetchMenu: async (menuID: number) => {
			set({ loading: true, error: null });
			try {
				const response = await menuBtnService.getMenuBtns(menuID);
				set({
					menuBtnData: response,
					loading: false,
				});
			} catch (err) {
				set({ error: (err as Error).message, loading: false });
			}
		},
		updateOrCreateMenu: async (value: MenuBtn) => {
			set({ loading: true, error: null });
			try {
				let newData: MenuBtn;

				if (value.id) {
					// 更新操作
					await menuBtnService.updateMenuBtn(value.id, value);
					newData = { ...value }; // 保持原有的数据，包括 id
				} else {
					// 创建操作
					const response = await menuBtnService.createMenuBtn(value);
					newData = { ...value, id: response.id }; // 将生成的 id 添加到 newData 中
				}

				// 更新状态，将 newData 添加到 menuBtnData 中或更新现有数据
				set((state) => {
					const existingIndex = state.menuBtnData.findIndex((item) => item.id === newData.id);
					if (existingIndex !== -1) {
						// 更新已有数据
						const updatedData = [...state.menuBtnData];
						updatedData[existingIndex] = newData;
						return { menuBtnData: updatedData, loading: false };
					}
					// 添加新数据
					return {
						menuBtnData: [...state.menuBtnData, newData],
						loading: false,
					};
				});
			} catch (err) {
				console.error(err);
				set({ error: (err as Error).message, loading: false });
			}
		},
		deleteMenu: async (id: number) => {
			try {
				await menuBtnService.deleteMenuBtn(id);
				// delete data
				set((state) => ({
					menuBtnData: state.menuBtnData.filter((item) => item.id !== id),
				}));
			} catch (err) {
				console.error(err);
			}
		},
	},
}));

export const useMenuBtnActions = () => useMenuBtnStore((state) => state.actions);
export const useMenuBtn = () => useMenuBtnStore((state) => state.menuBtnData);
export const useMenuBtnLoading = () => useMenuBtnStore((state) => state.loading);
export const useMenuBtnError = () => useMenuBtnStore((state) => state.error);
