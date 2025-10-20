package model

import (
	"time"
)

// UserRegisteredEvent 用户注册事件
type UserRegisteredEvent struct {
	ID           string
	UserID       string
	Username     string
	Email        string
	RegisteredAt time.Time
}

// EventID 事件ID
func (e *UserRegisteredEvent) EventID() string {
	return e.ID
}

// EventType 事件类型
func (e *UserRegisteredEvent) EventType() string {
	return UserRegisteredEventType
}

// Timestamp 事件时间戳
func (e *UserRegisteredEvent) Timestamp() time.Time {
	return e.RegisteredAt
}

// Payload 事件载荷
func (e *UserRegisteredEvent) Payload() interface{} {
	return map[string]interface{}{
		"userID":       e.UserID,
		"username":     e.Username,
		"email":        e.Email,
		"registeredAt": e.RegisteredAt,
	}
}
