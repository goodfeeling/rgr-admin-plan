import axios, { type AxiosRequestConfig, type AxiosError, type AxiosResponse } from "axios";

import { t } from "@/locales/i18n";
import userStore from "@/store/userStore";
import { Modal } from "antd";
import { toast } from "sonner";
import type { Result } from "#/api";
import { MessageType, PagePath, ResultEnum, StorageEnum } from "#/enum";
import userService from "./services/userService";

// create axios instance
const axiosInstance = axios.create({
	baseURL: import.meta.env.VITE_APP_BASE_API,
	timeout: 50000,
	headers: { "Content-Type": "application/json;charset=utf-8" },
});

// is refreshing
let isRefreshing = false;
// add a flag to track whether the logout modal has been shown
let isLogoutModalShown = false;

let failedQueue: Array<{
	resolve: (value?: any) => void;
	reject: (error?: any) => void;
}> = [];

const processQueue = (error: any, token: string | null = null) => {
	for (const prom of failedQueue) {
		if (error) {
			prom.reject(error);
		} else {
			prom.resolve(token);
		}
	}
	failedQueue = [];
};

// request interceptor
axiosInstance.interceptors.request.use(
	(config) => {
		const { userToken } = userStore.getState();
		if (userToken?.accessToken) {
			config.headers.Authorization = `Bearer ${userToken?.accessToken}`;
		}
		return config;
	},
	(error) => {
		return Promise.reject(error);
	},
);

// response interceptor
axiosInstance.interceptors.response.use(
	(res: AxiosResponse<Result>) => {
		// check is download file
		if (res.config.responseType === "blob" || res.data instanceof Blob) {
			// no handle return blob data
			return res;
		}
		if (!res.data) throw new Error(t("sys.api.apiRequestFailed"));
		const { status = 0, data, message = "" } = res.data;
		const hasSuccess = data && Reflect.has(res.data, "status") && status === ResultEnum.SUCCESS;

		if (hasSuccess) {
			return data;
		}
		throw new Error(message || t("sys.api.apiRequestFailed"));
	},
	async (error: AxiosError<Result>) => {
		const originalRequest = error.config as AxiosRequestConfig & {
			_retry?: boolean;
		};

		error.response?.status === 401 && handleAuthError(error);

		// checkout is 401 error and not a retry request
		if (error.response?.status === 401 && !originalRequest._retry) {
			if (isRefreshing) {
				// if already refreshing, add the request to the queue
				return new Promise((resolve, reject) => {
					failedQueue.push({ resolve, reject });
				})
					.then((token) => {
						originalRequest.headers = originalRequest.headers || {};
						originalRequest.headers.Authorization = `Bearer ${token}`;
						return axiosInstance(originalRequest);
					})
					.catch((err) => Promise.reject(err));
			}

			originalRequest._retry = true;
			isRefreshing = true;
			const { userToken, actions } = userStore.getState();

			if (!userToken?.refreshToken) {
				const errMsg = new Error("Invalid refresh");
				clearUserTokenToLoginPage(errMsg.message);
				return Promise.reject(errMsg);
			}

			try {
				const response = await userService.refreshToken(userToken?.refreshToken);

				const { jwtAccessToken, jwtRefreshToken, expirationAccessDateTime, expirationRefreshDateTime } =
					response.security;
				actions.setUserToken({
					accessToken: jwtAccessToken,
					refreshToken: jwtRefreshToken,
					expirationAccessDateTime,
					expirationRefreshDateTime,
				});
				processQueue(null, jwtAccessToken);
				return axiosInstance(originalRequest);
			} catch (err) {
				console.log(err);
				processQueue(err);
				clearUserTokenToLoginPage(err);
				return Promise.reject(new Error("Token refresh failed"));
			} finally {
				isRefreshing = false;
			}
		} else {
			const { response, message } = error || {};
			const newError = new Error(response?.data?.error || message || t("sys.api.errorMessage"));
			toast.error(newError.message, {
				position: "top-center",
			});
			return Promise.reject(newError);
		}
	},
);

// handler Error
function handleAuthError(error: AxiosError<Result<any>, any>) {
	const errorMessage = error.response?.data.error;

	const isTokenReplaced =
		errorMessage === "Token has been replaced" ||
		errorMessage === "refresh token has been replaced" ||
		errorMessage === "refresh token has been revoked";

	const isTokenInvalid =
		errorMessage === "Invalid token" || errorMessage === "Token expired" || errorMessage === "Token is expired";

	const isRefreshRequest = error.request.responseURL?.includes("/v1/auth/access-token");
	console.log(isTokenInvalid, isRefreshRequest);

	if (isTokenReplaced || (isTokenInvalid && isRefreshRequest)) {
		clearUserTokenToLoginPage(errorMessage);
		return Promise.reject(new Error(errorMessage));
	}
}

// clear user token
export function clearUserTokenToLoginPage(message: string | undefined) {
	// if already shown, do nothing
	if (isLogoutModalShown) {
		return;
	}
	// set true to show logout modal
	isLogoutModalShown = true;

	// clear localStorage in user store
	const storageKeys = [StorageEnum.STSToken, StorageEnum.UserStore, StorageEnum.Menu];
	for (const key of storageKeys) {
		localStorage.removeItem(key);
	}

	const { title, content } = (MessageType as Record<string, { title: string; content: string }>)[
		message as keyof typeof MessageType
	] || {
		title: "系统错误提示",
		content: message || "系统错误，请联系管理员",
	};

	Modal.warning({
		title,
		content,
		okText: "重新登录",
		centered: true,
		onOk: () => {
			window.location.replace(`#${PagePath.Login}`);
			window.location.reload();
		},
		// close modal set false
		afterClose: () => {
			isLogoutModalShown = false;
		},
	});
}

class APIClient {
	get<T = any>(config: AxiosRequestConfig): Promise<T> {
		return this.request({ ...config, method: "GET" });
	}

	post<T = any>(config: AxiosRequestConfig): Promise<T> {
		return this.request({ ...config, method: "POST" });
	}

	put<T = any>(config: AxiosRequestConfig): Promise<T> {
		return this.request({ ...config, method: "PUT" });
	}

	delete<T = any>(config: AxiosRequestConfig): Promise<T> {
		return this.request({ ...config, method: "DELETE" });
	}
	async download(config: AxiosRequestConfig): Promise<Blob> {
		const response = await axiosInstance({
			...config,
			responseType: "blob",
		});

		return response.data;
	}

	request<T = any>(config: AxiosRequestConfig): Promise<T> {
		return new Promise((resolve, reject) => {
			axiosInstance
				.request<any, AxiosResponse<Result>>(config)
				.then((res: AxiosResponse<Result>) => {
					resolve(res as unknown as Promise<T>);
				})
				.catch((e: Error | AxiosError) => {
					reject(e);
				});
		});
	}
}
export default new APIClient();
