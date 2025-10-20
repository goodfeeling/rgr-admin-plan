package model

import (
	"time"
)

// ApplicationEvent 应用事件接口
type ApplicationEvent interface {
	EventID() string
	EventType() string
	Timestamp() time.Time
	Payload() interface{}
}

// EventHandler 事件处理器接口
type EventHandler interface {
	Handle(event ApplicationEvent) error
}
