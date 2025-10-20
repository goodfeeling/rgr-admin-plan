import menuParameterService from "@/api/services/menuParameterService";
import type { MenuParameter } from "@/types/entity";
import { create } from "zustand";

interface MenuParameterState {
	menuParameterData: MenuParameter[];
	loading: boolean;
	error: string | null;
	actions: {
		fetchMenu: (MenuId: number) => Promise<void>;
		updateOrCreateMenu: (data: MenuParameter) => Promise<void>;
		deleteMenu: (id: number) => Promise<void>;
	};
}
export const useMenuParameterActions = () => useMenuParameterBtnStore((state) => state.actions);
export const useMenuParameter = () => useMenuParameterBtnStore((state) => state.menuParameterData);
export const useMenuParameterLoading = () => useMenuParameterBtnStore((state) => state.loading);
export const useMenuParameterError = () => useMenuParameterBtnStore((state) => state.error);

const useMenuParameterBtnStore = create<MenuParameterState>()((set) => ({
	menuParameterData: [],
	loading: true,
	error: null,

	actions: {
		fetchMenu: async (MenuId: number) => {
			set({ loading: true, error: null });
			try {
				const response = await menuParameterService.getMenuParameters(MenuId);
				set({
					menuParameterData: response,
					loading: false,
				});
			} catch (err) {
				set({ error: (err as Error).message, loading: false });
			}
		},
		updateOrCreateMenu: async (value: MenuParameter) => {
			set({ loading: true, error: null });
			try {
				let newData: MenuParameter;

				if (value.id) {
					// 更新操作
					await menuParameterService.updateMenuParameter(value.id, value);
					newData = { ...value }; // 保持原有的数据，包括 id
				} else {
					// 创建操作
					const response = await menuParameterService.createMenuParameter(value);
					newData = { ...value, id: response.id }; // 将生成的 id 添加到 newData 中
				}

				// 更新状态，将 newData 添加到 menuParameterData 中或更新现有数据
				set((state) => {
					const existingIndex = state.menuParameterData.findIndex((item) => item.id === newData.id);
					if (existingIndex !== -1) {
						// 更新已有数据
						const updatedData = [...state.menuParameterData];
						updatedData[existingIndex] = newData;
						return { menuParameterData: updatedData, loading: false };
					}
					// 添加新数据
					return {
						menuParameterData: [...state.menuParameterData, newData],
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
				await menuParameterService.deleteMenuParameter(id);
				// delete data
				set((state) => ({
					menuParameterData: state.menuParameterData.filter((item) => item.id !== id),
				}));
			} catch (err) {
				console.error(err);
			}
		},
	},
}));
