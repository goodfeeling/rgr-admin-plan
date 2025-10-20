import apiClient from "../apiClient";

import type { FileInfo, STSToken } from "#/entity";

export class UploadService {
	/**
	 * 单文件上传
	 */
	singleUpload() {
		return apiClient.get<FileInfo[]>({ url: `${UploadService.Client.Single}` });
	}

	/**
	 * 多文件上传
	 */
	multipleUpload() {
		return apiClient.get<FileInfo[]>({
			url: `${UploadService.Client.Multiple}`,
		});
	}

	/**
	 * 获取阿里云STS令牌
	 */
	getSTSToken() {
		return apiClient.get<STSToken>({
			url: UploadService.Client.AliYUnSTSToken,
		});
	}

	/**
	 * 刷新阿里云STS令牌
	 * @param refreshToken 刷新令牌
	 */
	refreshSTSToken(refreshToken: string) {
		return apiClient.get<STSToken>({
			url: `${UploadService.Client.RefreshSTSToken}?refresh_token=${refreshToken}`,
		});
	}
}

export namespace UploadService {
	export enum Client {
		Single = "/upload/single",
		Multiple = "/upload/multiple",
		AliYUnSTSToken = "/upload/sts-token",
		RefreshSTSToken = "/upload/refresh-sts",
	}
}

export default new UploadService();
