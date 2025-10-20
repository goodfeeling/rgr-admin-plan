package websocket

import (
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocketManager WebSocket管理器
type WebSocketManager struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	handlers   []WebSocketHandler // 存储多个处理器
	mutex      sync.RWMutex       // 保护handlers的并发访问
}

// WebSocketHandler WebSocket处理器接口
type WebSocketHandler interface {
	OnConnect(conn *websocket.Conn)
	OnMessage(conn *websocket.Conn, message []byte)
	OnDisconnect(conn *websocket.Conn)
}

// NewWebSocketManager 创建新的WebSocket管理器
func NewWebSocketManager() *WebSocketManager {
	manager := &WebSocketManager{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		handlers:   make([]WebSocketHandler, 0),
	}

	// 启动管理器
	go manager.Run()
	return manager
}

// AddHandler 添加处理器
func (wsm *WebSocketManager) AddHandler(handler WebSocketHandler) {
	wsm.mutex.Lock()
	defer wsm.mutex.Unlock()
	wsm.handlers = append(wsm.handlers, handler)
}

// RemoveHandler 移除处理器
func (wsm *WebSocketManager) RemoveHandler(handler WebSocketHandler) {
	wsm.mutex.Lock()
	defer wsm.mutex.Unlock()

	for i, h := range wsm.handlers {
		if h == handler {
			// 从切片中移除
			wsm.handlers = append(wsm.handlers[:i], wsm.handlers[i+1:]...)
			break
		}
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Upgrade 升级HTTP连接到WebSocket
func (wsm *WebSocketManager) Upgrade(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// HandleConnection 处理WebSocket连接
func (wsm *WebSocketManager) HandleConnection(c *gin.Context) {
	conn, err := wsm.Upgrade(c)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	// 注册连接
	wsm.register <- conn
	defer func() { wsm.unregister <- conn }()

	// 获取当前所有处理器
	wsm.mutex.RLock()
	handlers := make([]WebSocketHandler, len(wsm.handlers))
	copy(handlers, wsm.handlers)
	wsm.mutex.RUnlock()

	// 调用所有处理器的连接回调
	for _, handler := range handlers {
		handler.OnConnect(conn)
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
		handler.OnDisconnect(conn)
	}
}

// Run 运行WebSocket管理器
func (wsm *WebSocketManager) Run() {
	for {
		select {
		case conn := <-wsm.register:
			wsm.clients[conn] = true
			log.Println("Client connected. Total clients:", len(wsm.clients))

		case conn := <-wsm.unregister:
			if _, ok := wsm.clients[conn]; ok {
				delete(wsm.clients, conn)
				conn.Close()
				log.Println("Client disconnected. Total clients:", len(wsm.clients))
			}

		case message := <-wsm.broadcast:
			for conn := range wsm.clients {
				err := conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					delete(wsm.clients, conn)
					conn.Close()
				}
			}
		}
	}
}

// BroadcastMessage 广播消息给所有连接的客户端
func (wsm *WebSocketManager) BroadcastMessage(message []byte) {
	wsm.broadcast <- message
}

// SendMessageToClient 发送消息给特定客户端
func (wsm *WebSocketManager) SendMessageToClient(conn *websocket.Conn, message []byte) error {
	return conn.WriteMessage(websocket.TextMessage, message)
}

// GetClientsCount 获取当前连接数
func (wsm *WebSocketManager) GetClientsCount() int {
	return len(wsm.clients)
}
