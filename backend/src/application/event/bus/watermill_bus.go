// application/event/bus/watermill_bus.go
package bus

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/gbrayhan/microservices-go/src/application/event/model"
)

// WatermillEventBus 基于Watermill的事件总线
type WatermillEventBus struct {
	publisher    message.Publisher
	subscriber   message.Subscriber
	logger       watermill.LoggerAdapter
	handlerMutex sync.RWMutex
	handlers     map[string][]model.EventHandler // 跟踪每个事件类型的所有处理器
}

// NewWatermillEventBus 创建Watermill事件总线
func NewWatermillEventBus() (*WatermillEventBus, error) {
	logger := watermill.NewStdLogger(false, false)

	pubSub := gochannel.NewGoChannel(
		gochannel.Config{},
		logger,
	)

	return &WatermillEventBus{
		publisher:  pubSub,
		subscriber: pubSub,
		logger:     logger,
		handlers:   make(map[string][]model.EventHandler),
	}, nil
}

// Subscribe 订阅事件
func (eb *WatermillEventBus) Subscribe(eventType string, handler model.EventHandler) error {
	eb.handlerMutex.Lock()
	eb.handlers[eventType] = append(eb.handlers[eventType], handler)
	eb.handlerMutex.Unlock()

	messages, err := eb.subscriber.Subscribe(context.Background(), eventType)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			var eventModel map[string]interface{}
			if err := json.Unmarshal(msg.Payload, &eventModel); err != nil {
				eb.logger.Error("Cannot unmarshal event", err, nil)
				msg.Nack()
				continue
			}

			event := &GenericApplicationEvent{
				Data: eventModel,
				Type: eventType,
			}

			// 执行所有处理器
			eb.handlerMutex.RLock()
			handlers := eb.handlers[eventType]
			eb.handlerMutex.RUnlock()

			for _, h := range handlers {
				if err := h.Handle(event); err != nil {
					eb.logger.Error("Error handling event", err, nil)
					msg.Nack()
					continue
				}
			}

			msg.Ack()
		}
	}()

	return nil
}

// Publish 发布事件
func (eb *WatermillEventBus) Publish(ctx context.Context, event model.ApplicationEvent) error {
	payload, err := json.Marshal(event.Payload())
	if err != nil {
		return err
	}

	msg := message.NewMessage(event.EventID(), payload)
	return eb.publisher.Publish(event.EventType(), msg)
}

// Unsubscribe 取消订阅指定事件类型
func (eb *WatermillEventBus) Unsubscribe(eventType string, handler model.EventHandler) error {
	eb.handlerMutex.Lock()
	defer eb.handlerMutex.Unlock()

	handlers, exists := eb.handlers[eventType]
	if !exists {
		return nil
	}

	for i, h := range handlers {
		if h == handler {
			eb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	return nil
}

// GenericApplicationEvent 通用应用事件实现
type GenericApplicationEvent struct {
	Data map[string]interface{}
	Type string
	ID   string
	Time int64
}

func (e *GenericApplicationEvent) EventID() string {
	if e.ID != "" {
		return e.ID
	}
	id, ok := e.Data["id"].(string)
	if !ok {
		return ""
	}
	return id
}

func (e *GenericApplicationEvent) EventType() string {
	return e.Type
}

func (e *GenericApplicationEvent) Timestamp() time.Time {
	if e.Time != 0 {
		return time.Unix(e.Time, 0)
	}
	return time.Now()
}

func (e *GenericApplicationEvent) Payload() interface{} {
	return e.Data
}
