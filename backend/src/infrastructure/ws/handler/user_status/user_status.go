package user_status

import (
	ws "github.com/gbrayhan/microservices-go/src/infrastructure/lib/websocket"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
)

type UserStatusHandler struct {
	Logger         *logger.Logger
	sessionManager *ws.SessionManager
}

func NewUserStatusHandler(sessionManager *ws.SessionManager, loggerInstance *logger.Logger) *UserStatusHandler {
	return &UserStatusHandler{
		sessionManager: sessionManager,
		Logger:         loggerInstance,
	}
}

func (ch *UserStatusHandler) OnConnect(conn *websocket.Conn) {
	ch.Logger.Info("Log handler: Client connected")
	// 基础连接处理
}

func (ch *UserStatusHandler) OnConnectWithContext(conn *websocket.Conn, ctx *ws.WebSocketContext) {
	ch.Logger.Info("Log handler: Client OnConnectWithContext")
	// 从上下文获取用户信息
	userID := getUserIDFromContext(ctx.Context)
	deviceID := ctx.Context.Query("deviceId")
	ch.Logger.Info("Log handler: Client OnConnectWithContext:", zap.Int64("user_id", userID), zap.String("device_id", deviceID))
	// 注册用户会话
	ch.sessionManager.AddSession(userID, deviceID, conn)
}

func (ch *UserStatusHandler) OnMessage(conn *websocket.Conn, message []byte) {
	ch.Logger.Info("Log handler: Received message: %s", zap.String("Message", string(message)))
	// 处理心跳等消息
	// 可以实现心跳机制保持会话活跃
}

func (ch *UserStatusHandler) OnDisconnect(conn *websocket.Conn) {
	ch.Logger.Info("Log handler: Client disconnected")
	// 连接断开时清理会话

}
func (ch *UserStatusHandler) OnDisconnectWithContext(conn *websocket.Conn, ctx *ws.WebSocketContext) {
	ch.Logger.Info("Log handler: Client OnDisconnectWithContext")
	// 从上下文获取用户信息
	userID := getUserIDFromContext(ctx.Context)
	deviceID := ctx.Context.Query("deviceId")
	ch.Logger.Info("Log handler: Client OnConnectWithContext:", zap.Int64("user_id", userID), zap.String("device_id", deviceID))
	ch.sessionManager.RemoveSession(userID, deviceID)
}

// 辅助函数
func getUserIDFromContext(ctx *gin.Context) int64 {
	// 根据您的认证实现获取用户ID
	if userID, exists := ctx.Get("user_id"); exists {
		if id, ok := userID.(int); ok {
			return int64(id)
		}
	}
	return 0
}
