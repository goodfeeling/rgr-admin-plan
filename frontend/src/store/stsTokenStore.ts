import uploadService from "@/api/services/uploadService";
import type { STSToken } from "@/types/entity";
import { StorageEnum } from "@/types/enum";
import OSS from "ali-oss";
import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

interface STSTokenState {
	stsToken: STSToken | null;
	loading: boolean;
	error: string | null;
	actions: {
		fetchSTSToken: () => Promise<void>;
		refreshSTSToken: () => Promise<void>;
		uploadToOSS: (
			file: File,
			customName?: string,
		) => Promise<{
			success?: boolean;
			url?: string;
			name?: string;
			error?: string;
		}>;
	};
}

const useSTSTokenStore = create<STSTokenState>()(
	persist(
		(set, get) => ({
			stsToken: null,
			loading: true,
			error: null,

			actions: {
				fetchSTSToken: async () => {
					set({ loading: true, error: null });
					try {
						const response = await uploadService.getSTSToken();
						set({
							stsToken: response,
							loading: false,
						});
					} catch (err) {
						set({ error: (err as Error).message, loading: false });
					}
				},

				refreshSTSToken: async () => {
					set({ loading: true, error: null });
					try {
						const currentStsToken = get().stsToken;
						if (!currentStsToken || !currentStsToken.refresh_token) {
							throw new Error("No refresh token available");
						}

						// 调用刷新接口，这里假设后端支持使用 refresh_token 刷新 STS Token
						const response = await uploadService.refreshSTSToken(currentStsToken.refresh_token);
						set({
							stsToken: response,
							loading: false,
						});
					} catch (err) {
						set({ error: (err as Error).message, loading: false });
						// 如果刷新失败，可能需要重新获取新的 STS Token
						await get().actions.fetchSTSToken();
					}
				},

				uploadToOSS: async (file: File, customName?: string) => {
					try {
						const stsToken = get().stsToken;

						// 检查 STS Token 是否存在
						if (!stsToken) {
							throw new Error("STS Token is not available");
						}
						// 配置 OSS 客户端
						const client = new OSS({
							region: stsToken.region,
							accessKeyId: stsToken.access_key_id,
							accessKeySecret: stsToken.access_key_secret,
							stsToken: stsToken.security_token,
							bucket: stsToken.bucket_name,
						});

						// 生成文件名
						const fileName = customName || `${Date.now()}-${file.name}`;

						// 上传文件
						const result = await client.put(fileName, file, {
							// 添加进度回调
							progress: (p: any) => {
								console.log("上传进度:", p);
							},
						});

						return {
							success: true,
							url: result.url,
							name: fileName,
						};
					} catch (error: any) {
						console.error("OSS上传失败:", error);

						// 提供更详细的错误信息
						let errorMessage = "未知错误";

						// 检查是否是CORS错误
						if (error.code === "RequestError" && error.status === -1) {
							errorMessage = "CORS配置错误或网络连接问题，请检查OSS的跨域设置";
						} else if (error.code) {
							switch (error.code) {
								case "AccessDenied":
									errorMessage = "访问被拒绝，请检查STS Token权限";
									break;
								case "NoSuchBucket":
									errorMessage = "Bucket不存在，请检查配置";
									break;
								default:
									errorMessage = error.message || error.code;
							}
						} else if (error.message) {
							errorMessage = error.message;
						}

						return {
							success: false,
							error: errorMessage,
						};
					}
				},
			},
		}),
		{
			name: StorageEnum.STSToken,
			storage: createJSONStorage(() => localStorage),
			partialize: (state) => ({ [StorageEnum.STSToken]: state.stsToken }),
		},
	),
);

export const useSTSTokenActions = () => useSTSTokenStore((state) => state.actions);
export const useSTSToken = () => useSTSTokenStore((state) => state.stsToken);

export const useSTSTokenLoading = () => useSTSTokenStore((state) => state.loading);
export const useSTSTokenError = () => useSTSTokenStore((state) => state.error);
