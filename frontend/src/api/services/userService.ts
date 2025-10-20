import apiClient from "../apiClient";

import type { PageList, PasswordEditReq, UpdateUser, UserInfo } from "#/entity";

export class UserService {
	/**
	 * 用户登录
	 * @param data 登录信息
	 */
	signin(data: UserService.SignInReq) {
		return apiClient.post<UserService.SignInRes>({
			url: UserService.Client.SignIn,
			data,
		});
	}

	/**
	 * 用户注册
	 * @param data 注册信息
	 */
	signup(data: UserService.SignUpReq) {
		return apiClient.post<UserService.SignInRes>({
			url: UserService.Client.SignUp,
			data,
		});
	}

	/**
	 * 用户登出
	 */
	logout() {
		return apiClient.get({ url: UserService.Client.Logout });
	}

	/**
	 * 刷新token
	 * @param refreshToken 刷新token
	 */
	refreshToken(refreshToken: string) {
		return apiClient.post({
			url: UserService.Client.Refresh,
			data: {
				refreshToken,
			},
		});
	}

	/**
	 * 根据ID查找用户
	 * @param id 用户ID
	 */
	findById(id: string) {
		return apiClient.get<UserInfo[]>({
			url: `${UserService.Client.User}/${id}`,
		});
	}

	/**
	 * 更新用户信息
	 * @param id 用户ID
	 * @param userInfo 用户信息
	 */
	updateUser(id: number, userInfo: UpdateUser) {
		return apiClient.put<UserInfo>({
			url: `${UserService.Client.User}/${id}`,
			data: userInfo,
		});
	}

	/**
	 * 创建用户
	 * @param userInfo 用户信息
	 */
	createUser(userInfo: UserInfo) {
		return apiClient.post<UserInfo>({
			url: `${UserService.Client.User}`,
			data: userInfo,
		});
	}

	/**
	 * 搜索用户分页列表
	 * @param searchStr 搜索字符串
	 */
	searchPageList(searchStr: string) {
		return apiClient.get<PageList<UserInfo>>({
			url: `${UserService.Client.SearchUser}?${searchStr}`,
		});
	}

	/**
	 * 删除用户
	 * @param id 用户ID
	 */
	deleteUser(id: number) {
		return apiClient.delete<string>({
			url: `${UserService.Client.User}/${id}`,
		});
	}

	/**
	 * 绑定角色
	 * @param userId 用户ID
	 * @param roleIds 角色ID列表
	 */
	bindRole(userId: number, roleIds: string[]) {
		return apiClient.post<boolean>({
			url: `${UserService.Client.User}/${userId}/role`,
			data: {
				roleIds,
			},
		});
	}

	/**
	 * 重置密码
	 * @param id 用户ID
	 */
	resetPassword(id: number) {
		return apiClient.post<boolean>({
			url: `${UserService.Client.User}/${id}/reset-password`,
		});
	}

	/**
	 * 编辑密码
	 * @param id 用户ID
	 * @param updateInfo 更新信息
	 */
	editPassword(id: number, updateInfo: PasswordEditReq) {
		return apiClient.post<boolean>({
			url: `${UserService.Client.User}/${id}/edit-password`,
			data: updateInfo,
		});
	}

	/**
	 * 更改密码
	 * @param resetInfo 重置信息
	 * @param resetToken 重置token
	 */
	changePassword(resetInfo: UserService.PasswordResetReq, resetToken: string) {
		return apiClient.post<boolean>({
			url: `${UserService.Client.User}/change-password?token=${resetToken}`,
			data: resetInfo,
		});
	}

	/**
	 * 切换角色
	 * @param roleId 角色ID
	 */
	switchRole(roleId: number) {
		return apiClient.post<UserService.SignInRes>({
			url: `${UserService.Client.SwitchRole}?role_id=${roleId}`,
		});
	}
}

export namespace UserService {
	export enum Client {
		SignIn = "/auth/signin",
		SignUp = "/auth/signup",
		Logout = "/auth/logout",
		Refresh = "/auth/access-token",
		SwitchRole = "/auth/switch-role",

		User = "/user",
		SearchUser = "/user/search",
		UserStatusWs = "/ws/user/status",
	}
	export interface SignInReq {
		user_name: string;
		password: string;
		captcha_answer: string;
		captcha_id: string;
	}

	export interface SignUpReq extends SignInReq {
		email: string;
	}

	export type SignInRes = {
		security: {
			expirationAccessDateTime: string;
			expirationRefreshDateTime: string;
			jwtAccessToken: string;
			jwtRefreshToken: string;
		};
		userinfo: UserInfo;
	};

	export type PasswordResetReq = {
		email: string;
		new_password: string;
		confirm_password: string;
	};
}

export default new UserService();
