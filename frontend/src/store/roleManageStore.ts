import roleService from "@/api/services/roleService";
import type { Role } from "@/types/entity";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { create } from "zustand";

interface RoleManageState {
	data: Role[];
}

const useRoleManageStore = create<RoleManageState>()(() => ({
	data: [],
}));

// 更新
export const useUpdateOrCreateRoleMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async (data: Role) => {
			if (data.id) {
				const { parent_id = 0, name = "", label = "", order = 0, description = "", status = 0 } = data;
				await roleService.updateRole(data.id, {
					parent_id,
					name,
					label,
					order,
					description,
					status,
				});
				return { ...data };
			}
			// 创建
			const response = await roleService.createRole(data);
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
export const useRemoveRoleMutation = () => {
	const queryClient = useQueryClient();

	return useMutation({
		mutationFn: async (id: number) => {
			await roleService.deleteRole(id);
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

export const useRoleQuery = () => {
	return useQuery({
		queryKey: ["apiManageList"],
		queryFn: () => {
			return roleService.getRoles();
		},
	});
};

export const useRoleManage = () => useRoleManageStore((state) => state.data);
