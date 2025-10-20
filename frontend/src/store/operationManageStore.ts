import operationService from "@/api/services/operationService";
import type { Operation, PageList, TableParams } from "@/types/entity";
import { getRandomUserParams, toURLSearchParams } from "@/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { create } from "zustand";

interface OperationManageState {
	data: PageList<Operation>;
	condition: TableParams;
	actions: {
		setCondition: (tableParams: TableParams) => void;
	};
}

const useOperationManageStore = create<OperationManageState>()((set) => ({
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

// 删除
export const useRemoveOperationMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await operationService.deleteOperation(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["operationManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};
// 批量删除
export const useBatchRemoveOperationMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (selectedRowKeys: number[]) => {
			await operationService.deleteBatch(selectedRowKeys);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["operationManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

export const useOperationQuery = () => {
	const tableParams = useOperationManageStore.getState().condition;
	return useQuery({
		queryKey: [
			"operationManageList",
			tableParams.pagination?.current,
			tableParams.pagination?.pageSize,
			tableParams.sortField,
			tableParams.sortOrder,
			tableParams.searchParams,
			tableParams.filters,
		],
		queryFn: () => {
			const params = toURLSearchParams(
				getRandomUserParams(tableParams, (result, searchParams) => {
					if (searchParams) {
						if (searchParams.method) {
							result.method_like = searchParams.method;
						}
						if (searchParams.status) {
							result.status_match = searchParams.status;
						}
						if (searchParams.path) {
							result.path_like = searchParams.path;
						}
					}
				}),
			);
			return operationService.searchPageList(params.toString());
		},
	});
};

export const useOperationManage = () => useOperationManageStore((state) => state.data);

export const useOperationManageCondition = () => useOperationManageStore((state) => state.condition);
export const useOperationActions = () => useOperationManageStore((state) => state.actions);
