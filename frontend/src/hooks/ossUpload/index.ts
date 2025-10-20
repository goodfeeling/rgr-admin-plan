import { useSTSToken, useSTSTokenActions } from "@/store/stsTokenStore";

export const useOssUpload = () => {
	const stsToken = useSTSToken();
	const actions = useSTSTokenActions();
	return {
		uploadFile: async (file: File, customName?: string) => {
			let currentStsToken = stsToken;

			let retryCount = 0;
			const maxRetries = 2;

			while ((!currentStsToken || !currentStsToken.access_key_id) && retryCount < maxRetries) {
				if (!currentStsToken || !currentStsToken.access_key_id) {
					// 如果没有STS Token，获取一个新的
					await actions.fetchSTSToken();
				}

				// 等待一小段时间再检查
				await new Promise((resolve) => setTimeout(resolve, 100));

				// 重新获取状态
				currentStsToken = useSTSToken();
				retryCount++;
			}

			if (!currentStsToken) {
				// 如果没有STS Token，获取一个新的
				await actions.fetchSTSToken();
				// 获取新的Token后重新获取状态
				currentStsToken = useSTSToken();
			}

			if (currentStsToken?.expiration) {
				const expirationDate = new Date(currentStsToken.expiration);
				const now = new Date();
				const timeDiff = expirationDate.getTime() - now.getTime();
				const fiveMinutesInMs = 5 * 60 * 1000;

				// 如果距离过期不足5分钟，刷新Token
				if (timeDiff < fiveMinutesInMs) {
					try {
						await actions.refreshSTSToken();
						// 刷新后重新获取Token
						currentStsToken = useSTSToken();
					} catch (error) {
						console.error("刷新STS Token失败:", error);
						// 刷新失败时获取新的Token
						await actions.fetchSTSToken();
						// 获取新的Token后重新获取状态
						currentStsToken = useSTSToken();
					}
				}
			}

			// 确保有有效的Token再执行上传
			if (currentStsToken) {
				return await actions.uploadToOSS(file, customName);
			}
			// 如果仍然没有有效的Token，返回失败结果
			return {
				success: false,
				error: "无法获取有效的STS Token",
			};
		},
	};
};
