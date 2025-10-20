import dictionaryService from "@/api/services/dictionaryService";
import type { Dictionary, PageList, TableParams } from "@/types/entity";
import { getRandomUserParams, toURLSearchParams } from "@/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { create } from "zustand";

interface DictionaryManageState {
	data: PageList<Dictionary>;
	condition: TableParams;
	actions: {
		setCondition: (tableParams: TableParams) => void;
	};
}

const useDictionaryManageStore = create<DictionaryManageState>()((set) => ({
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
export const useUpdateOrCreateDictionaryMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async (data: Dictionary) => {
			if (data.id) {
				await dictionaryService.updateDictionary(data.id, data);
				return { ...data };
			}
			// 创建
			const response = await dictionaryService.createDictionary(data);
			return { ...data, id: response.id };
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["dictionaryManageList"] });
		},
		onError: (err) => {
			console.error("Update or create API failed:", err);
		},
	});
};

// 删除
export const useRemoveDictionaryMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await dictionaryService.deleteDictionary(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["dictionaryManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};
// 批量删除
export const useBatchRemoveDictionaryMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (selectedRowKeys: number[]) => {
			await dictionaryService.deleteBatch(selectedRowKeys);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["dictionaryManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

export const useDictionaryQuery = () => {
	const tableParams = useDictionaryManageStore.getState().condition;
	return useQuery({
		queryKey: [
			"dictionaryManageList",
			tableParams.pagination?.current,
			tableParams.pagination?.pageSize,
			tableParams.sortField,
			tableParams.sortOrder,
			tableParams.searchParams,
			tableParams.filters,
		],
		queryFn: () => {
			const params = toURLSearchParams(getRandomUserParams(tableParams));
			return dictionaryService.searchPageList(params.toString());
		},
	});
};

export const useDictionaryManage = () => useDictionaryManageStore((state) => state.data);

export const useDictionaryManageCondition = () => useDictionaryManageStore((state) => state.condition);
export const useDictionaryActions = () => useDictionaryManageStore((state) => state.actions);
