// src/utils/deviceId.ts
import { StorageEnum } from "#/enum";
import { getItem, setItem } from "./storage";

export const getOrCreateDeviceId = (): string => {
	// 从localStorage中获取已存在的设备ID
	let deviceId = getItem<string>(StorageEnum.DeviceId);

	// 如果不存在，则生成一个新的设备ID
	if (!deviceId) {
		deviceId = generateComplexDeviceId();
		setItem(StorageEnum.DeviceId, deviceId);
	}

	return deviceId;
};

// 生成更复杂的设备ID
const generateComplexDeviceId = (): string => {
	// 结合多种信息生成设备ID
	const navigatorInfo = [
		navigator.userAgent,
		navigator.platform,
		navigator.language,
		new Date().getTimezoneOffset(),
	].join("|");

	// 使用简单的hash函数（实际项目中可以使用更复杂的hash算法）
	let hash = 0;
	for (let i = 0; i < navigatorInfo.length; i++) {
		const char = navigatorInfo.charCodeAt(i);
		hash = (hash << 5) - hash + char;
		hash = hash & hash; // 转换为32位整数
	}

	return `device_${Math.abs(hash).toString(36)}${Date.now().toString(36)}`;
};
