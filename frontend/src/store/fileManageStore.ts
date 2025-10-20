import fileService from "@/api/services/fileService";
import type { FileInfo, PageList, TableParams } from "@/types/entity";
import { getRandomUserParams, toURLSearchParams } from "@/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { create } from "zustand";

interface FileInfoManageState {
	data: PageList<FileInfo>;
	condition: TableParams;
	actions: {
		setCondition: (tableParams: TableParams) => void;
	};
}

const useFileInfoManageStore = create<FileInfoManageState>()((set) => ({
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
export const useUpdateOrCreateFileInfoMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async (data: FileInfo) => {
			if (data.id) {
				await fileService.updateFileInfo(data.id, data);
				return { ...data };
			}
			// 创建
			const response = await fileService.createFileInfo(data);
			return { ...data, id: response.id };
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["fileManageList"] });
		},
		onError: (err) => {
			console.error("Update or create API failed:", err);
		},
	});
};

// 删除
export const useRemoveFileInfoMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await fileService.deleteFileInfo(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["fileManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};
// 批量删除
export const useBatchRemoveFileInfoMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (selectedRowKeys: number[]) => {
			await fileService.deleteBatch(selectedRowKeys);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["fileManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

export const useFileInfoQuery = () => {
	const tableParams = useFileInfoManageStore.getState().condition;
	return useQuery({
		queryKey: [
			"fileManageList",
			tableParams.pagination?.current,
			tableParams.pagination?.pageSize,
			tableParams.sortField,
			tableParams.sortOrder,
			tableParams.searchParams?.storage_engine,
			tableParams.searchParams?.file_origin_name,
			tableParams.filters,
		],
		queryFn: () => {
			const params = toURLSearchParams(
				getRandomUserParams(tableParams, (result, searchParams) => {
					if (searchParams) {
						if (searchParams.file_origin_name) {
							result.file_origin_name_like = searchParams.file_origin_name;
						}
						if (searchParams.storage_engine) {
							result.storage_engine_like = searchParams.storage_engine;
						}
					}
				}),
			);

			return fileService.searchPageList(params.toString());
		},
	});
};

export const useFileInfoManage = () => useFileInfoManageStore((state) => state.data);

export const useFileInfoManageCondition = () => useFileInfoManageStore((state) => state.condition);
export const useFileInfoActions = () => useFileInfoManageStore((state) => state.actions);
