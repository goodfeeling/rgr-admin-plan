import { usePathname } from "@/routes/hooks";
import { useMenu } from "@/store/useMenuStore";
import { useUserInfo } from "@/store/userStore";
import { Button } from "@/ui/button";
import type React from "react";
import { getPathnames, groupCheck } from "./common";

interface PermissionButtonProps {
	permissionString: string;
	children: React.ReactNode;
	fallback?: React.ReactNode;
	className?: string;
	[key: string]: any;
}

const PermissionButton: React.FC<PermissionButtonProps> = ({
	permissionString,
	children,
	fallback = null,
	className,
	...restProps
}) => {
	const userInfo = useUserInfo();
	const menuGroup = useMenu();
	const pathnames = getPathnames(usePathname());
	const menuData = groupCheck(menuGroup, pathnames);

	// 检查是否具有按钮权限
	const hasPermission = () => {
		// 如果是超级管理员，拥有所有权限
		if (userInfo.id === 1) {
			return true;
		}

		// 检查当前角色是否拥有该按钮权限
		return menuData !== null && !!menuData.btn_slice?.includes(permissionString);
	};

	// 如果有权限，渲染按钮，否则渲染 fallback 或 null
	if (hasPermission()) {
		return (
			<Button className={className} {...restProps}>
				{children}
			</Button>
		);
	}

	return fallback ? <>{fallback}</> : null;
};

export default PermissionButton;
