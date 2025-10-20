import menuGroupService from "@/api/services/menuGroupService";
import type { MenuGroup, PageList, TableParams } from "@/types/entity";
import { getRandomUserParams, toURLSearchParams } from "@/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { create } from "zustand";

interface MenuGroupManageState {
	data: PageList<MenuGroup>;
	condition: TableParams;
	actions: {
		setCondition: (tableParams: TableParams) => void;
	};
}

const useMenuGroupManageStore = create<MenuGroupManageState>()((set) => ({
	data: {
		list: [],
		total: 0,
		page: 1,
		page_size: 10,
		filters: undefined,
		total_page: 1,
	},
	condition: {
		pagination: {
			current: 1,
			pageSize: 10,
			total: 0,
		},
		sortField: "id",
		sortOrder: "descend",
	},
	actions: {
		setCondition: (condition: TableParams) => {
			set({ condition });
		},
	},
}));

// 更新
export const useUpdateOrCreateMenuGroupMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async (data: MenuGroup) => {
			if (data.id) {
				await menuGroupService.updateMenuGroup(data.id, data);
				return { ...data };
			}
			// 创建
			const response = await menuGroupService.createMenuGroup(data);
			return { ...data, id: response.id };
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["menuGroupManageList"] });
		},
		onError: (err) => {
			console.error("Update or create API failed:", err);
		},
	});
};

// 删除
export const useRemoveMenuGroupMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await menuGroupService.deleteMenuGroup(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["menuGroupManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

export const useMenuGroupQuery = () => {
	const tableParams = useMenuGroupManageStore.getState().condition;
	return useQuery({
		queryKey: [
			"menuGroupManageList",
			tableParams.pagination?.current,
			tableParams.pagination?.pageSize,
			tableParams.sortField,
			tableParams.sortOrder,
			tableParams.searchParams,
			tableParams.filters,
		],
		queryFn: () => {
			const params = toURLSearchParams(getRandomUserParams(tableParams));
			return menuGroupService.searchPageList(params.toString());
		},
	});
};

export const useMenuGroupManage = () => useMenuGroupManageStore((state) => state.data);

export const useMenuGroupManageCondition = () => useMenuGroupManageStore((state) => state.condition);
export const useMenuGroupActions = () => useMenuGroupManageStore((state) => state.actions);
