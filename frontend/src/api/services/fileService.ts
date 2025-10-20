import apiClient from "../apiClient";

import type { FileInfo, PageList } from "#/entity";

export class FileService {
	/**
	 * 获取文件列表
	 */
	getFileInfos() {
		return apiClient.get<FileInfo[]>({
			url: `${FileService.Client.FileInfo}`,
		});
	}

	/**
	 * 更新文件信息
	 * @param id 文件ID
	 * @param apiInfo 文件信息
	 */
	updateFileInfo(id: number, apiInfo: FileInfo) {
		return apiClient.put<FileInfo>({
			url: `${FileService.Client.FileInfo}/${id}`,
			data: apiInfo,
		});
	}

	/**
	 * 创建文件信息
	 * @param apiInfo 文件信息
	 */
	createFileInfo(apiInfo: FileInfo) {
		return apiClient.post<FileInfo>({
			url: `${FileService.Client.FileInfo}`,
			data: apiInfo,
		});
	}

	/**
	 * 搜索文件分页列表
	 * @param searchStr 搜索字符串
	 */
	searchPageList(searchStr: string) {
		return apiClient.get<PageList<FileInfo>>({
			url: `${FileService.Client.SearchFileInfo}?${searchStr}`,
		});
	}

	/**
	 * 删除文件
	 * @param id 文件ID
	 */
	deleteFileInfo(id: number) {
		return apiClient.delete<string>({
			url: `${FileService.Client.FileInfo}/${id}`,
		});
	}

	/**
	 * 批量删除文件
	 * @param ids 文件ID数组
	 */
	deleteBatch(ids: number[]) {
		return apiClient.post<number>({
			url: `${FileService.Client.DeleteBatch}`,
			data: { ids },
		});
	}
}

export namespace FileService {
	export enum Client {
		FileInfo = "/file",
		SearchFileInfo = "/file/search",
		GroupsFileInfo = "/file/groups",
		DeleteBatch = "/file/batch",
	}
}

export default new FileService();
