import apiClient from "../apiClient";

import type { ConfigResponse } from "#/entity";

export class ConfigService {
	/**
	 * 获取所有配置
	 */
	getConfigs() {
		return apiClient.get<ConfigResponse>({
			url: `${ConfigService.Client.Config}`,
		});
	}

	/**
	 * 更新配置
	 * @param dataInfo 配置数据
	 * @param module 模块名称
	 */
	updateConfig(dataInfo: { [key: string]: string }, module: string) {
		return apiClient.put<{ [key: string]: string }>({
			url: `${ConfigService.Client.Config}/${module}`,
			data: dataInfo,
		});
	}

	/**
	 * 根据模块获取配置
	 * @param module 模块名称
	 */
	getConfigByModule(module: string) {
		return apiClient.get<{ [key: string]: string }>({
			url: `${ConfigService.Client.ConfigModule}/${module}`,
		});
	}

	/**
	 * 获取站点配置
	 */
	getConfigBySite() {
		return apiClient.get<{ [key: string]: string }>({
			url: `${ConfigService.Client.Config}/site`,
		});
	}
}

export namespace ConfigService {
	export enum Client {
		Config = "/config",
		ConfigModule = "/config/module",
	}
}

export default new ConfigService();
