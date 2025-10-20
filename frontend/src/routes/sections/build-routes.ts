// routes/buildRoutes.ts
import type { Menu, MenuTreeUserGroup } from "@/types/entity";
import React from "react";
import { Navigate, type RouteObject } from "react-router";
import { lazyLoad } from "./lazy-wrapper";
const modules = import.meta.glob("/src/pages/**/*.tsx");
export function buildRoutes(menuData: Menu[]): RouteObject[] {
	return menuData.map((item) => {
		const path = `/src/pages/${item.component}.tsx`;
		const LazyComponent = React.lazy(() => modules[path]() as Promise<{ default: React.ComponentType<any> }>);

		const route: RouteObject = {
			path: item.path,
		};

		if (item.children && item.children.length > 0) {
			// 递归生成子路由
			const children = buildRoutes(item.children);

			// 添加默认跳转（跳到第一个子菜单的 path）
			const firstChildPath = item.children[0]?.path;
			if (firstChildPath) {
				children.unshift({
					index: true,
					element: React.createElement(Navigate, {
						to: firstChildPath,
						replace: true,
					}),
				});
			}

			route.children = children;
		} else {
			route.element = lazyLoad(LazyComponent, {
				keepAlive: item.keep_alive === 1,
				keepAliveName: item.path,
			});
		}

		return route;
	});
}

export function convertMenuTreeUserGroupToMenus(data: MenuTreeUserGroup[]): Menu[] {
	return data.map((group, index) => {
		return {
			// 映射字段
			name: group.name || "-",
			path: group.path || "-",
			title: group.name || "-",

			// 必填的其它字段，可以给默认值或从 group.items 中推断
			id: index + 1, // 如果没有 id 可用索引或 UUID
			menu_level: 1, // 默认设置为第一级
			parent_id: 0,
			hidden: false,
			component: `${group.path}/index`, // 如果有默认组件名，也可以传
			sort: 0,
			keep_alive: 0,
			icon: "", // 你可以从 group.name 推图标
			menu_group_id: 0,
			created_at: new Date().toISOString(),
			updated_at: new Date().toISOString(),
			level: [1],

			// 如果你想保留 group.items 作为子菜单，也可以映射进去
			children: group.items || [],
		};
	});
}
