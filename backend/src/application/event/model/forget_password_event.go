package model

import (
	"time"
)

// ForgetPasswordEvent 发送邮件事件
type ForgetPasswordEvent struct {
	ID           string
	To           string
	Subject      string
	Body         string
	RegisteredAt time.Time
}

// EventID 事件ID
func (e *ForgetPasswordEvent) EventID() string {
	return e.ID
}

// EventType 事件类型
func (e *ForgetPasswordEvent) EventType() string {
	return ForgetPasswordEventType
}

// Timestamp 事件时间戳
func (e *ForgetPasswordEvent) Timestamp() time.Time {
	return e.RegisteredAt
}

// Payload 事件载荷
func (e *ForgetPasswordEvent) Payload() interface{} {
	return map[string]interface{}{
		"to":      e.To,
		"subject": e.Subject,
		"body":    e.Body,
	}
}
