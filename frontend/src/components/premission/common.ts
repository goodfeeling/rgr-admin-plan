import type { Menu, MenuTreeUserGroup } from "@/types/entity";

const checkChildren = (menu: Menu, pathnames: string[], level: number): Menu | null => {
	if (!pathnames[level] || !menu) {
		return null;
	}
	if (menu.path === `/${pathnames.join("/")}`) {
		return menu;
	}
	const currentPath = menu.path.replace("/", "").split("/");
	if (currentPath[currentPath.length - 1] === pathnames[level]) {
		const isLastPath = level === pathnames.length - 1;
		if (!menu.children || menu.children.length <= 0) {
			return isLastPath ? menu : null;
		}

		if (menu.children && menu.children.length > 0) {
			for (const child of menu.children) {
				const result = checkChildren(child, pathnames, level + 1);
				if (result !== null) {
					return result;
				}
			}
		}
		return null;
	}
	return null;
};

// 顶层遍历
export const groupCheck = (menuGroup: MenuTreeUserGroup[], pathnames: string[]): Menu | null => {
	if (!menuGroup || !pathnames || pathnames.length === 0) {
		return null;
	}
	for (const item of menuGroup) {
		if (item.path === pathnames[0]) {
			for (const menuItem of item.items) {
				const result = checkChildren(menuItem, pathnames, 1);
				if (result !== null) {
					return result;
				}
			}
		}
	}
	return null;
};

export const getPathnames = (pathname: string): string[] => {
	const rawPathname = pathname.replace("/", "");
	if (rawPathname === "") {
		return [];
	}
	return rawPathname.split("/");
};
