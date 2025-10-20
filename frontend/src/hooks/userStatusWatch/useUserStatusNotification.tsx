import { clearUserTokenToLoginPage } from "@/api/apiClient";
import { UserService } from "@/api/services/userService";
import { useSharedWebSocket } from "@/hooks/sharedWebSocket";
import { useUserToken } from "@/store/userStore";
import { getOrCreateDeviceId } from "@/utils/deviceId";
import { useEffect } from "react";

const useUserStatusNotification = () => {
	const { accessToken } = useUserToken();
	const deviceId = getOrCreateDeviceId(); // 获取设备ID
	// 使用已有的WebSocket hook监听用户状态
	const { message } = useSharedWebSocket(accessToken ? `${UserService.Client.UserStatusWs}?deviceId=${deviceId}` : "");
	useEffect(() => {
		if (!accessToken) return;
		if (message) {
			// 根据服务器发送的消息处理用户状态变化
			if (message.type === "FORCE_LOGOUT") {
				clearUserTokenToLoginPage(message.message);
			}
			// 可以根据实际消息格式进行调整
		}
	}, [message, accessToken]);
};

export default useUserStatusNotification;
