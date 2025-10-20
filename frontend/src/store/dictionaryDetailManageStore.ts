import dictionaryDetailService from "@/api/services/dictionaryDetailService";
import type { DictionaryDetail, PageList, TableParams } from "@/types/entity";
import { getRandomUserParams, toURLSearchParams } from "@/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { create } from "zustand";

interface DictionaryDetailManageState {
	data: PageList<DictionaryDetail>;
	condition: TableParams;
	actions: {
		setCondition: (tableParams: TableParams) => void;
	};
}

const useDictionaryDetailManageStore = create<DictionaryDetailManageState>()((set) => ({
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
		sortField: "sort",
		sortOrder: "descend",
	},
	actions: {
		setCondition: (condition: TableParams) => {
			set({ condition });
		},
	},
}));

// 更新
export const useUpdateOrCreateDictionaryDetailMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async (data: DictionaryDetail) => {
			if (data.id) {
				await dictionaryDetailService.updateDictionaryDetail(data.id, data);
				return { ...data };
			}
			// 创建
			const response = await dictionaryDetailService.createDictionaryDetail(data);
			return { ...data, id: response.id };
		},
		onSuccess: () => {
			queryClient.invalidateQueries({
				queryKey: ["dictionaryDetailManageList"],
			});
		},
		onError: (err) => {
			console.error("Update or create API failed:", err);
		},
	});
};

// 删除
export const useRemoveDictionaryDetailMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await dictionaryDetailService.deleteDictionaryDetail(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({
				queryKey: ["dictionaryDetailManageList"],
			});
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};
// 批量删除
export const useBatchRemoveDictionaryDetailMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (selectedRowKeys: number[]) => {
			await dictionaryDetailService.deleteBatch(selectedRowKeys);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({
				queryKey: ["dictionaryDetailManageList"],
			});
		},
		onError: (err) => {
			console.error("Delete  failed:", err);
		},
	});
};

export const useDictionaryDetailQuery = () => {
	const tableParams = useDictionaryDetailManageStore.getState().condition;

	return useQuery({
		queryKey: [
			"dictionaryDetailManageList",
			tableParams.pagination?.current,
			tableParams.pagination?.pageSize,
			tableParams.sortField,
			tableParams.sortOrder,
			tableParams.searchParams?.selectedDictId,
			tableParams.filters,
		],
		queryFn: () => {
			const params = toURLSearchParams(
				getRandomUserParams(tableParams, (result, searchParams) => {
					if (searchParams.selectedDictId) {
						result.selectedDictId_match = searchParams.selectedDictId;
					}
				}),
			);
			return dictionaryDetailService.searchPageList(params.toString());
		},
	});
};

export const useDictionaryDetailManage = () => useDictionaryDetailManageStore((state) => state.data);

export const useDictionaryDetailManageCondition = () => useDictionaryDetailManageStore((state) => state.condition);
export const useDictionaryDetailActions = () => useDictionaryDetailManageStore((state) => state.actions);
