import configService from "@/api/services/configService";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";

// 更新
export const useUpdateOrCreateConfigMutation = () => {
	const queryClient = useQueryClient();
	return useMutation({
		mutationFn: async ({
			data,
			module,
		}: {
			data: { [key: string]: any };
			module: string;
		}) => {
			await configService.updateConfig(data, module);
			return { ...data };
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["configManageList"] });
		},
		onError: (err) => {
			console.error("Update or create API failed:", err);
		},
	});
};

export const useConfigQuery = () => {
	return useQuery({
		queryKey: ["configManageList"],
		queryFn: () => {
			return configService.getConfigs();
		},
	});
};
