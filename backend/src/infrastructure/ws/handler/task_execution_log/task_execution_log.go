package task_execution_log

import (
	"encoding/json"
	"sync"

	domainTaskExecutionLog "github.com/gbrayhan/microservices-go/src/domain/sys/task_execution_log"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	ws "github.com/gbrayhan/microservices-go/src/infrastructure/lib/websocket"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// 扩展 WebSocketContext 来存储订阅信息
type ClientSubscription struct {
	Conn    *websocket.Conn
	TaskIDs map[int]bool
}

// LogHandler 日志处理器
type LogHandler struct {
	Logger                  *logger.Logger
	taskExecutionLogService domainTaskExecutionLog.ITaskExecutionLogService
	subscriptions           map[*websocket.Conn]*ClientSubscription
	subscriptionMutex       sync.RWMutex
}

func NewLogHandler(
	taskExecutionLogService domainTaskExecutionLog.ITaskExecutionLogService,
	loggerInstance *logger.Logger,
) *LogHandler {
	return &LogHandler{
		Logger:                  loggerInstance,
		taskExecutionLogService: taskExecutionLogService,
		subscriptions:           make(map[*websocket.Conn]*ClientSubscription),
		subscriptionMutex:       sync.RWMutex{},
	}
}

// 实现扩展接口
func (ch *LogHandler) OnConnectWithContext(conn *websocket.Conn, ctx *ws.WebSocketContext) {
	ch.Logger.Info("Log handler: Client connected")
	ch.subscriptionMutex.Lock()
	defer ch.subscriptionMutex.Unlock()
	// 初始化客户端的订阅信息
	ch.subscriptions[conn] = &ClientSubscription{
		Conn:    conn,
		TaskIDs: make(map[int]bool),
	}
	ch.Logger.Info("Added new subscription", zap.Int("total_subscriptions", len(ch.subscriptions)))

}

func (ch *LogHandler) OnConnect(conn *websocket.Conn) {
	ch.Logger.Info("Log handler: Client connected")

}

func (ch *LogHandler) OnMessage(conn *websocket.Conn, message []byte) {
	ch.Logger.Info("Log handler: Received message: %s", zap.String("Message", string(message)))

	// 解析传入的 JSON 字符串
	var request struct {
		TaskID int `json:"taskId"`
		Limit  int `json:"limit"`
	}
	if err := json.Unmarshal(message, &request); err != nil {
		ch.sendError(conn, "Invalid JSON format: "+err.Error())
		return
	}
	// 添加订阅
	ch.subscriptionMutex.Lock()
	if subscription, exists := ch.subscriptions[conn]; exists {
		subscription.TaskIDs[request.TaskID] = true
	}

	ch.subscriptionMutex.Unlock()
	ch.Logger.Info("Client subscribed to task", zap.Int("TaskID", request.TaskID))
	// 设置默认值
	if request.Limit <= 0 {
		request.Limit = 100
	}
	// 调用 GetByTaskID 获取日志数据
	result, err := ch.taskExecutionLogService.GetByTaskID(uint(request.TaskID), request.Limit)
	if err != nil {
		ch.sendError(conn, "Failed to fetch logs: "+err.Error())
		return
	}
	ch.Logger.Info("Log handler: Fetched logs: %v", zap.Any("Logs", len(*result)))

	// 发送结果
	ch.sendData(conn, "logs", result)
}

// 向特定 taskID 的订阅者推送日志
func (ch *LogHandler) NotifyLogToTaskSubscribers(taskID int, logData interface{}) {
	ch.Logger.Info("Log handler: Notifying subscribers for task ID: ", zap.Int("TaskID", taskID))
	ch.subscriptionMutex.RLock()
	defer ch.subscriptionMutex.RUnlock()

	response := map[string]interface{}{
		"type":   "log_update",
		"taskId": taskID,
		"data":   logData,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		ch.Logger.Error("Failed to marshal log notification", zap.Error(err))
		return
	}
	// 遍历所有订阅者，向订阅了该 taskID 的客户端发送消息
	for _, subscription := range ch.subscriptions {
		if subscription.TaskIDs[taskID] {
			err := subscription.Conn.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				ch.Logger.Error("Failed to send log notification", zap.Error(err))
			}
		}
	}
}

func (ch *LogHandler) OnDisconnect(conn *websocket.Conn) {
	ch.Logger.Info("Log handler: Client disconnected")

	// 清理订阅信息
	ch.subscriptionMutex.Lock()
	delete(ch.subscriptions, conn)
	ch.subscriptionMutex.Unlock()
}
func (ch *LogHandler) OnDisconnectWithContext(conn *websocket.Conn, ctx *ws.WebSocketContext) {
	ch.Logger.Info("Log handler: Client disconnected")
	// 连接断开时清理会话

}

// 辅助方法
func (ch *LogHandler) sendData(conn *websocket.Conn, msgType string, data interface{}) {
	response := map[string]interface{}{
		"type": msgType,
		"data": data,
	}

	jsonData, err := json.Marshal(response)
	if err != nil {
		ch.sendError(conn, "Failed to marshal response")
		return
	}

	conn.WriteMessage(websocket.TextMessage, jsonData)
}

func (ch *LogHandler) sendError(conn *websocket.Conn, message string) {
	errorResponse := map[string]interface{}{
		"type":  "error",
		"error": message,
	}

	jsonData, _ := json.Marshal(errorResponse)
	conn.WriteMessage(websocket.TextMessage, jsonData)
}
