import { useMutation } from "@tanstack/react-query";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

import userService, { type UserService } from "@/api/services/userService";

import { toast } from "sonner";
import type { PasswordEditReq, UserInfo, UserToken } from "#/entity";
import { StorageEnum } from "#/enum";
import { useMenuActions } from "./useMenuStore";

type UserStore = {
	userInfo: Partial<UserInfo>;
	userToken: UserToken;
	// 使用 actions 命名空间来存放所有的 action
	actions: {
		setUserInfo: (userInfo: UserInfo) => void;
		setUserToken: (token: UserToken) => void;
		switchRole: (roleId: number) => Promise<void>;
		clearUserInfoAndToken: () => void;
		passwordEdit: (id: number, editInfo: PasswordEditReq) => Promise<void>;
	};
};

const useUserStore = create<UserStore>()(
	persist(
		(set) => ({
			userInfo: {},
			userToken: {},
			actions: {
				setUserInfo: (userInfo) => {
					set({ userInfo });
				},
				setUserToken: (userToken) => {
					set({ userToken });
				},
				switchRole: async (roleId: number): Promise<void> => {
					return new Promise<void>((resolve, reject) => {
						(async () => {
							try {
								const response = await userService.switchRole(roleId);
								const { security } = response;
								const { jwtAccessToken, jwtRefreshToken, expirationAccessDateTime, expirationRefreshDateTime } =
									security;
								set({
									userToken: {
										accessToken: jwtAccessToken,
										refreshToken: jwtRefreshToken,
										expirationAccessDateTime,
										expirationRefreshDateTime,
									},
									userInfo: response.userinfo,
								});
								resolve();
							} catch (error) {
								console.log(error);
								reject(error);
							}
						})();
					});
				},
				passwordEdit: async (id: number, editInfo: PasswordEditReq) => {
					return new Promise<void>((resolve, reject) => {
						(async () => {
							try {
								await userService.editPassword(id, editInfo);

								resolve();
							} catch (error) {
								console.log(error);
								reject(error);
							}
						})();
					});
				},
				clearUserInfoAndToken() {
					set({ userInfo: {}, userToken: {} });
				},
			},
		}),
		{
			name: "userStore", // name of the item in the storage (must be unique)
			storage: createJSONStorage(() => localStorage), // (optional) by default, 'localStorage' is used
			partialize: (state) => ({
				[StorageEnum.UserInfo]: state.userInfo,
				[StorageEnum.UserToken]: state.userToken,
			}),
		},
	),
);

export const useUserInfo = () => useUserStore((state) => state.userInfo);
export const useUserToken = () => useUserStore((state) => state.userToken);
export const useUserActions = () => useUserStore((state) => state.actions);

export const useSignIn = () => {
	const { setUserToken, setUserInfo } = useUserActions();
	const menuActions = useMenuActions();

	const signInMutation = useMutation({
		mutationFn: userService.signin,
	});

	const signIn = async (data: UserService.SignInReq) => {
		try {
			const res = await signInMutation.mutateAsync(data);
			const { userinfo: userInfo, security } = res;
			const { jwtAccessToken, jwtRefreshToken, expirationAccessDateTime, expirationRefreshDateTime } = security;
			setUserToken({
				accessToken: jwtAccessToken,
				refreshToken: jwtRefreshToken,
				expirationAccessDateTime,
				expirationRefreshDateTime,
			});
			setUserInfo(userInfo);
			// get user menu
			await menuActions.fetchMenu();
			toast.success("Sign in success!", {
				closeButton: true,
			});
		} catch (err) {
			toast.error(err.message, {
				position: "top-center",
			});
		}
	};
	return signIn;
};
export default useUserStore;
