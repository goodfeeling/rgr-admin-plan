import { usePathname } from "@/routes/hooks";
import { useMenu } from "@/store/useMenuStore";
import { useUserInfo } from "@/store/userStore";
import type React from "react";
import { useCallback, useEffect } from "react";
import { useNavigate } from "react-router";
import { getPathnames, groupCheck } from "./common";

interface PermissionGuardProps {
	children: React.ReactNode;
	fallback?: React.ReactNode;
}

const PermissionGuard: React.FC<PermissionGuardProps> = ({ children, fallback = null }) => {
	const userInfo = useUserInfo();
	const menuGroup = useMenu();
	const pathnames = getPathnames(usePathname());
	const navigate = useNavigate();
	const menu = groupCheck(menuGroup, pathnames);
	// 检查是否具有权限
	const hasPermission = useCallback(() => {
		// 登录进来会到"/",pathnames === []，所以跳过权限校验
		if (pathnames.length === 0) {
			return true;
		}
		// 如果是超级管理员，拥有所有权限
		if (userInfo.id === 1) {
			return true;
		}

		// 检查当前角色是否拥有该路由权限
		return menu !== null;
	}, [menu, pathnames, userInfo.id]);

	useEffect(() => {
		if (!hasPermission()) {
			navigate("/403");
		}
	}, [hasPermission, navigate]);

	// 如果有权限，渲染子组件，否则渲染 fallback 或 null
	return hasPermission() ? <>{children}</> : <>{fallback}</>;
};

export default PermissionGuard;
