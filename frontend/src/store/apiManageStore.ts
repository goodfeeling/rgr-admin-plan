import apisService from "@/api/services/apisService";
import type { Api, PageList, TableParams } from "@/types/entity";
import { getRandomUserParams, toURLSearchParams } from "@/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";
import { create } from "zustand";

interface ApiManageState {
	data: PageList<Api>;
	condition: TableParams;
	actions: {
		setCondition: (tableParams: TableParams) => void;
	};
}

const useApiManageStore = create<ApiManageState>()((set) => ({
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
export const useUpdateOrCreateApiMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async (data: Api) => {
			if (data.id) {
				await apisService.updateApi(data.id, data);
				return { ...data };
			}
			// 创建
			const response = await apisService.createApi(data);
			return { ...data, id: response.id };
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["apiManageList"] });
		},
		onError: (err) => {
			console.error("Update or create API failed:", err);
		},
	});
};

// 删除
export const useRemoveApiMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await apisService.deleteApi(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["apiManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};
// 批量删除
export const useBatchRemoveApiMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (selectedRowKeys: number[]) => {
			await apisService.deleteBatch(selectedRowKeys);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["apiManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

// 同步数据
export const useSynchronizeApiMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async () => {
			return await apisService.synchronizeApi();
		},
		onSuccess: (response) => {
			queryClient.invalidateQueries({ queryKey: ["apiManageList"] });
			toast.success(`同步完成数量：${response.count} 条`);
		},
		onError: (err) => {
			console.error("Synchronize API failed:", err);
			toast.error(`同步失败: ${err.message}`);
		},
	});
};

export const useApiQuery = () => {
	const tableParams = useApiManageStore.getState().condition;
	return useQuery({
		queryKey: [
			"apiManageList",
			tableParams.pagination?.current,
			tableParams.pagination?.pageSize,
			tableParams.sortField,
			tableParams.sortOrder,
			tableParams.searchParams?.path,
			tableParams.searchParams?.description,
			tableParams.searchParams?.method,
			tableParams.searchParams?.api_group,
			tableParams.filters,
		],
		queryFn: () => {
			const params = toURLSearchParams(
				getRandomUserParams(tableParams, (result, searchParams) => {
					if (searchParams) {
						if (searchParams.path) {
							result.path_like = searchParams.path;
						}
						if (searchParams.description) {
							result.description_like = searchParams.description;
						}

						if (searchParams.method) {
							result.method_match = searchParams.method;
						}
						if (searchParams.api_group) {
							result.apiGroup_match = searchParams.api_group;
						}
					}
				}),
			);
			return apisService.searchPageList(params.toString());
		},
	});
};

export const useApiManage = () => useApiManageStore((state) => state.data);

export const useApiManageCondition = () => useApiManageStore((state) => state.condition);
export const useApiActions = () => useApiManageStore((state) => state.actions);
