import { getPathnames, groupCheck } from "@/components/premission/common";
import { useMenu } from "@/store/useMenuStore";
import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { useLocation } from "react-router";
import { useTabOperations } from "../hooks/use-tab-operations";
import type { KeepAliveTab, MultiTabsContextType } from "../types";
import { menuPathCache } from "./menu-path-cache";

const MultiTabsContext = createContext<MultiTabsContextType>({
	tabs: [],
	activeTabRoutePath: "",
	setTabs: () => {},
	closeTab: () => {},
	closeOthersTab: () => {},
	closeAll: () => {},
	closeLeft: () => {},
	closeRight: () => {},
	refreshTab: () => {},
});

export function MultiTabsProvider({ children }: { children: React.ReactNode }) {
	const [tabs, setTabs] = useState<KeepAliveTab[]>([]);
	const location = useLocation();
	const menuData = useMenu();
	const currentRouteKey = location.pathname;
	// 使用缓存查找菜单项
	const currentMenu = useMemo(() => {
		// 首先检查缓存中是否有结果
		if (menuPathCache.has(currentRouteKey)) {
			return menuPathCache.get(currentRouteKey);
		}

		// 如果缓存中没有，则执行查找
		const pathnames = getPathnames(currentRouteKey);
		const menu = groupCheck(menuData, pathnames);

		// 将结果存入缓存
		menuPathCache.set(currentRouteKey, menu);

		return menu;
	}, [currentRouteKey, menuData]);

	const activeTabRoutePath = useMemo(() => {
		return currentRouteKey;
	}, [currentRouteKey]);

	const operations = useTabOperations(tabs, setTabs, activeTabRoutePath);

	useEffect(() => {
		setTabs((prev) => {
			const filtered = prev.filter((item) => !item.hideTab);
			const isAlreadyExisted = filtered.some((item) => item.key === currentRouteKey);
			if (!isAlreadyExisted && currentRouteKey !== "/") {
				return [
					...filtered,
					{
						key: currentRouteKey,
						label: currentMenu?.title || currentRouteKey,
						hideTab: false,
						icon: currentMenu?.icon || "",
						children: null,
						params: {},
						timeStamp: new Date().getTime().toString(),
					},
				];
			}

			return filtered;
		});
	}, [currentRouteKey, currentMenu]);

	const contextValue = useMemo(
		() => ({
			tabs,
			activeTabRoutePath,
			setTabs,
			...operations,
		}),
		[tabs, activeTabRoutePath, operations],
	);

	return <MultiTabsContext.Provider value={contextValue}>{children}</MultiTabsContext.Provider>;
}

export function useMultiTabsContext() {
	return useContext(MultiTabsContext);
}
