package websocket

import (
	"log"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// RouteHandler 带路由信息的处理器
type RouteHandler struct {
	Route   string
	Handler WebSocketHandler
}

// WebSocketRouter WebSocket路由管理器
type WebSocketRouter struct {
	*WebSocketManager
	routes map[string][]WebSocketHandler
}

// WebSocketContext WebSocket上下文信息
type WebSocketContext struct {
	Context     *gin.Context
	QueryParams url.Values
	Route       string
}

// ExtendedWebSocketHandler 扩展的WebSocket处理器接口
type ExtendedWebSocketHandler interface {
	WebSocketHandler
	OnConnectWithContext(conn *websocket.Conn, ctx *WebSocketContext)
	OnDisconnectWithContext(conn *websocket.Conn, ctx *WebSocketContext)
}

// NewWebSocketRouter 创建新的WebSocket路由管理器
func NewWebSocketRouter() *WebSocketRouter {
	manager := NewWebSocketManager()
	return &WebSocketRouter{
		WebSocketManager: manager,
		routes:           make(map[string][]WebSocketHandler),
	}
}

// AddRoute 添加路由处理器
func (wsr *WebSocketRouter) AddRoute(route string, handler WebSocketHandler) {
	if _, exists := wsr.routes[route]; !exists {
		wsr.routes[route] = make([]WebSocketHandler, 0)
	}
	wsr.routes[route] = append(wsr.routes[route], handler)
}

// HandleConnectionWithRoute 处理带路由的WebSocket连接
func (wsr *WebSocketRouter) HandleConnectionWithRoute(c *gin.Context, route string) {
	conn, err := wsr.Upgrade(c)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 注册连接
	wsr.register <- conn
	defer func() { wsr.unregister <- conn }()

	// 获取指定路由的处理器
	handlers, exists := wsr.routes[route]
	if !exists {
		log.Printf("No handlers found for route: %s", route)
		return
	}

	// 创建扩展的处理器上下文，包含查询参数
	extendedContext := &WebSocketContext{
		Context:     c,
		QueryParams: c.Request.URL.Query(),
		Route:       route,
	}

	// 调用所有处理器的连接回调，传递上下文信息
	for _, handler := range handlers {
		// 如果处理器实现了扩展接口，则传递上下文
		if extendedHandler, ok := handler.(ExtendedWebSocketHandler); ok {
			extendedHandler.OnConnectWithContext(conn, extendedContext)
		} else {
			// 保持向后兼容
			handler.OnConnect(conn)
		}
	}

	// 处理消息循环
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		if messageType == websocket.TextMessage {
			// 调用所有处理器的消息回调
			for _, handler := range handlers {
				handler.OnMessage(conn, message)
			}
		}
	}

	// 调用所有处理器的断开连接回调
	for _, handler := range handlers {
		// 如果处理器实现了扩展接口，则传递上下文
		if extendedHandler, ok := handler.(ExtendedWebSocketHandler); ok {
			extendedHandler.OnDisconnectWithContext(conn, extendedContext)
		} else {
			// 保持向后兼容
			handler.OnDisconnect(conn)
		}
	}
}
