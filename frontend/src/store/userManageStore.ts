import userService from "@/api/services/userService";
import type { PageList, TableParams, UserInfo } from "@/types/entity";
import { getRandomUserParams, toURLSearchParams } from "@/utils";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

import { create } from "zustand";

interface UserManageState {
	data: PageList<UserInfo>;
	condition: TableParams;
	actions: {
		setCondition: (tableParams: TableParams) => void;
	};
}

const useUserManageStore = create<UserManageState>()((set) => ({
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
export const useUpdateOrCreateUserMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async (data: UserInfo) => {
			if (data.id) {
				await userService.updateUser(data.id, data);
				return { ...data };
			}
			// 创建
			const response = await userService.createUser(data);
			return { ...data, id: response.id };
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["userManageList"] });
		},
		onError: (err) => {
			console.error("Update or create API failed:", err);
		},
	});
};

// 删除
export const useRemoveUserMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await userService.deleteUser(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["userManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

// passwordEdit
export const usePasswordResetMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await userService.resetPassword(id);
		},
		onSuccess: () => {
			// 成功后使相关查询失效，触发重新获取
			queryClient.invalidateQueries({ queryKey: ["userManageList"] });
		},
		onError: (err) => {
			console.error("Delete API failed:", err);
		},
	});
};

export const useUserQuery = () => {
	const tableParams = useUserManageStore.getState().condition;
	return useQuery({
		queryKey: [
			"userManageList",
			tableParams.pagination?.current,
			tableParams.pagination?.pageSize,
			tableParams.sortField,
			tableParams.sortOrder,
			tableParams?.searchParams?.user_name,
			tableParams?.searchParams?.status,
			tableParams.filters,
		],
		queryFn: () => {
			const params = toURLSearchParams(
				getRandomUserParams(tableParams, (result, searchParams) => {
					if (searchParams) {
						if (searchParams.user_name) {
							result.userName_like = searchParams.user_name;
						}
						if (searchParams.status) {
							result.status_match = searchParams.status;
						}
					}
				}),
			);
			return userService.searchPageList(params.toString());
		},
	});
};

export const useUserManage = () => useUserManageStore((state) => state.data);

export const useUserManageCondition = () => useUserManageStore((state) => state.condition);
export const useUserManageActions = () => useUserManageStore((state) => state.actions);
