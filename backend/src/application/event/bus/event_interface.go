package bus

import (
	"context"

	"github.com/gbrayhan/microservices-go/src/application/event/model"
)

// EventBus 事件总线接口
type EventBus interface {
	Publish(ctx context.Context, event model.ApplicationEvent) error
	Subscribe(eventType string, handler model.EventHandler) error
	Unsubscribe(eventType string, handler model.EventHandler) error
}
