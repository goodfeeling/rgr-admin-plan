import { webSocketManager } from "@/utils/webSocketManager";
import { useCallback, useEffect, useState } from "react";

export const useSharedWebSocket = (uri: string) => {
	const url = `${import.meta.env.VITE_APP_WS_BASE_URL || "ws://127.0.0.1:8080"}${uri}`;
	const [_, setConnected] = useState(false);
	const [message, setMessage] = useState<any>(null);

	useEffect(() => {
		// url is empty, do not establish a connection
		if (!uri) {
			return;
		}

		const unsubscribe = webSocketManager.connect(url, (data) => {
			setMessage(data);
			setConnected(true);
		});

		return () => {
			unsubscribe();
			setConnected(false);
			setMessage(null);
		};
	}, [url, uri]);

	const sendMessage = useCallback(
		(data: string | ArrayBuffer | Blob | ArrayBufferView) => {
			// uri is empty string, do not send message
			if (!uri) {
				return;
			}
			webSocketManager.sendMessage(url, data);
		},
		[url, uri],
	);

	return {
		connected: uri ? webSocketManager.isConnected(url) : false,
		message,
		sendMessage,
	};
};
