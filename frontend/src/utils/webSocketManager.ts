// src/utils/webSocketManager.ts
import userStore from "@/store/userStore";

class WebSocketManager {
	private connections: Map<string, { ws: WebSocket; refCount: number; listeners: Set<(data: any) => void> }> =
		new Map();

	connect(url: string, onMessage: (data: any) => void): () => void {
		// 获取用户token
		const { userToken } = userStore.getState();
		const token = userToken?.accessToken;

		// 将token添加到URL参数中
		let fullUrl = url;
		if (token) {
			const separator = url.includes("?") ? "&" : "?";
			fullUrl = `${url}${separator}token=${encodeURIComponent(token)}`;
		}

		let connection = this.connections.get(fullUrl);

		if (!connection) {
			// 创建新连接
			const ws = new WebSocket(fullUrl);
			const listeners = new Set<(data: any) => void>();
			listeners.add(onMessage);

			ws.onopen = () => {
				console.log("WebSocket connected to", fullUrl);
			};

			ws.onmessage = (event) => {
				try {
					const data = JSON.parse(event.data);
					// 通知所有监听器
					for (const listener of listeners) {
						listener(data);
					}
				} catch (error) {
					console.error("Failed to parse WebSocket message:", error);
				}
			};

			ws.onerror = (error) => {
				console.error("WebSocket error:", error);
			};

			ws.onclose = () => {
				console.log("WebSocket disconnected from", fullUrl);
				this.connections.delete(fullUrl);
			};

			connection = { ws, refCount: 1, listeners };
			this.connections.set(fullUrl, connection);
		} else {
			// 增加引用计数
			connection.refCount++;
			connection.listeners.add(onMessage);
		}

		// 返回取消订阅函数
		return () => {
			if (connection) {
				connection.listeners.delete(onMessage);
				connection.refCount--;

				// 如果没有组件在使用，关闭连接
				if (connection.refCount <= 0) {
					connection.ws.close();
					this.connections.delete(fullUrl);
				}
			}
		};
	}

	sendMessage(url: string, data: string | ArrayBuffer | Blob | ArrayBufferView) {
		// 获取用户token
		const { userToken } = userStore.getState();
		const token = userToken?.accessToken;

		// 将token添加到URL参数中
		let fullUrl = url;
		if (token) {
			const separator = url.includes("?") ? "&" : "?";
			fullUrl = `${url}${separator}token=${encodeURIComponent(token)}`;
		}

		const connection = this.connections.get(fullUrl);
		if (connection && connection.ws.readyState === WebSocket.OPEN) {
			connection.ws.send(data);
		} else {
			console.warn("WebSocket is not connected to", fullUrl);
		}
	}

	isConnected(url: string): boolean {
		// 获取用户token
		const { userToken } = userStore.getState();
		const token = userToken?.accessToken;

		// 将token添加到URL参数中
		let fullUrl = url;
		if (token) {
			const separator = url.includes("?") ? "&" : "?";
			fullUrl = `${url}${separator}token=${encodeURIComponent(token)}`;
		}

		const connection = this.connections.get(fullUrl);
		return connection ? connection.ws.readyState === WebSocket.OPEN : false;
	}
}

export const webSocketManager = new WebSocketManager();
