import scheduleService from "@/api/services/scheduleService";
import type { PageList, ScheduledTask, TableParams } from "@/types/entity";
import { getRandomUserParams, toURLSearchParams } from "@/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { create } from "zustand";

interface ScheduledTaskManageState {
	data: PageList<ScheduledTask>;
	condition: TableParams;
	actions: {
		setCondition: (tableParams: TableParams) => void;
	};
}

const useScheduledTaskManageStore = create<ScheduledTaskManageState>()((set) => ({
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
export const useUpdateOrCreateScheduledTaskMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async (data: ScheduledTask) => {
			if (data.id) {
				await scheduleService.updateScheduledTask(data.id, data);
				return { ...data };
			}
			// 创建
			const response = await scheduleService.createScheduledTask(data);
			return { ...data, id: response.id };
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["scheduleManageList"] });
		},
		onError: (err) => {
			console.error("Update or create API failed:", err);
		},
	});
};

// 删除
export const useRemoveScheduledTaskMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await scheduleService.deleteScheduledTask(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["scheduleManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

// 批量删除
export const useBatchRemoveScheduledTaskMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (selectedRowKeys: number[]) => {
			await scheduleService.deleteBatch(selectedRowKeys);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["scheduleManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

// 启动任务
export const useEnableTaskMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await scheduleService.enableTask(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["scheduleManageList"] });
		},
		onError: (err) => {
			console.error("Enable API failed:", err);
		},
	});
};

// 关闭任务
export const useDisableTaskMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await scheduleService.disableTask(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["scheduleManageList"] });
		},
		onError: (err) => {
			console.error("Disable API failed:", err);
		},
	});
};

// 重启任务
export const useReloadTaskMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async () => {
			await scheduleService.reloadTask();
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["scheduleManageList"] });
		},
		onError: (err) => {
			console.error("Reload API failed:", err);
		},
	});
};

export const useScheduledTaskQuery = (options?: {
	enablePolling?: boolean;
}) => {
	const tableParams = useScheduledTaskManageStore.getState().condition;
	return useQuery({
		queryKey: [
			"scheduleManageList",
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
						if (searchParams.task_name) {
							result.taskName_like = searchParams.task_name;
						}

						if (searchParams.task_type) {
							result.taskType_like = searchParams.task_type;
						}
						if (searchParams.status) {
							result.status_match = searchParams.status;
						}
					}
				}),
			);
			return scheduleService.searchPageList(params.toString());
		},
		refetchInterval: options?.enablePolling ? 5000 : false,
	});
};

export const useScheduledTaskManage = () => useScheduledTaskManageStore((state) => state.data);

export const useScheduledTaskManageCondition = () => useScheduledTaskManageStore((state) => state.condition);
export const useScheduledTaskManegeActions = () => useScheduledTaskManageStore((state) => state.actions);
