import configService from "@/api/services/configService";
import { useQuery } from "@tanstack/react-query";

export function useMapBySystemConfig() {
	return useQuery({
		queryKey: ["sysConfig"],
		queryFn: () => configService.getConfigBySite(),
		staleTime: 1000 * 60 * 5,
		gcTime: 1000 * 60 * 30,
	});
}
