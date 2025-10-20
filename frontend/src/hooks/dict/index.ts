import dictionaryService from "@/api/services/dictionaryService";
import { useQuery } from "@tanstack/react-query";

type MapByTypeResult = { [key: string]: string };

// 添加使用React Query的缓存版本
export function useDictionaryByTypeWithCache(typeText: string) {
	return useQuery({
		queryKey: ["dictionary", typeText],
		queryFn: () => dictionaryService.getByType(typeText),
		staleTime: 1000 * 60 * 5,
		gcTime: 1000 * 60 * 30,
		select: (data) => data.details,
	});
}

export function useMapByTypeWithCache(typeText: string) {
	return useQuery({
		queryKey: ["dictionaryMap", typeText],
		queryFn: () => dictionaryService.getByType(typeText),
		staleTime: 1000 * 60 * 5,
		gcTime: 1000 * 60 * 30,
		select: (data) => {
			const result: MapByTypeResult = {};
			for (const item of data.details) {
				result[item.label] = item.value;
			}
			return result;
		},
	});
}
