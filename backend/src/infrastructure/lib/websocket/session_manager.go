// session_manager.go
package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type UserSession struct {
	UserID      int64
	DeviceID    string
	Conn        *websocket.Conn
	ConnectedAt time.Time
}

type SessionManager struct {
	sessions map[int64][]*UserSession // userID -> sessions
	mutex    sync.RWMutex
}

func NewSessionManager() *SessionManager {
	return &SessionManager{
		sessions: make(map[int64][]*UserSession),
	}
}
func (sm *SessionManager) AddSession(userID int64, deviceID string, conn *websocket.Conn) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	session := &UserSession{
		UserID:      userID,
		DeviceID:    deviceID,
		Conn:        conn,
		ConnectedAt: time.Now(),
	}

	sm.sessions[userID] = append(sm.sessions[userID], session)
}

func (sm *SessionManager) RemoveSession(userID int64, deviceID string) {
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	if sessions, exists := sm.sessions[userID]; exists {
		filtered := make([]*UserSession, 0)
		for _, session := range sessions {
			if session.DeviceID != deviceID {
				filtered = append(filtered, session)
			}
		}
		sm.sessions[userID] = filtered
	}
}

func (sm *SessionManager) NotifyOtherDevicesOffline(userID int64, currentDeviceID string) {
	sm.mutex.RLock()
	sessions, exists := sm.sessions[userID]
	sm.mutex.RUnlock()
	if !exists {
		return
	}
	for _, session := range sessions {
		// 通知除当前设备外的其他设备下线
		if session.DeviceID != currentDeviceID {
			message := map[string]interface{}{
				"type":      "FORCE_LOGOUT",
				"message":   "You have been logged in from another device",
				"timestamp": time.Now().Unix(),
			}

			// 异步发送，避免阻塞
			go func(conn *websocket.Conn) {
				conn.WriteJSON(message)
				conn.Close()
			}(session.Conn)
		}
	}

	// 清理已通知的会话
	sm.mutex.Lock()
	defer sm.mutex.Unlock()
	// 可以选择清理旧会话或保留当前会话
}
