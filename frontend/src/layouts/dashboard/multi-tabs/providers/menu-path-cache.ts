import type { Menu } from "@/types/entity";

// 创建一个可导出的缓存对象，供外部使用
export const menuPathCache = new Map<string, Menu | null>();

// 提供清除缓存的方法
export function clearMenuPathCache() {
	menuPathCache.clear();
}
