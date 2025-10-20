import menuService from "@/api/services/menuService";
import type { Menu } from "@/types/entity";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { create } from "zustand";

interface MenuManageState {
	data: Menu[];
}

const useMenuManageStore = create<MenuManageState>()(() => ({
	data: [],
}));

// 更新
export const useUpdateOrCreateMenuMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async (data: Menu) => {
			if (data.id) {
				await menuService.updateMenu(data.id, data);
				return { ...data };
			}
			// 创建
			const response = await menuService.createMenu(data);
			return { ...data, id: response.id };
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["menuManageList"] });
		},
		onError: (err) => {
			console.error("Update or create API failed:", err);
		},
	});
};

// 删除
export const useRemoveMenuMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await menuService.deleteMenu(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["menuManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

export const useMenuQuery = (selectedId: number) => {
	return useQuery({
		queryKey: ["menuManageList", selectedId],
		queryFn: () => {
			return menuService.getMenus(selectedId);
		},
		enabled: selectedId !== 0,
	});
};
export const useMenuManage = () => useMenuManageStore((state) => state.data);
