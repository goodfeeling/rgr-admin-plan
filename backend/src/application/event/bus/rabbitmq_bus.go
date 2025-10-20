package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gbrayhan/microservices-go/src/application/event/model"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/streadway/amqp"
)

// RabbitMQEventBus RabbitMQ事件总线实现
type RabbitMQEventBus struct {
	connection   *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
	queueName    string
	consumerTag  string
	handlers     map[string][]model.EventHandler
	handlerMutex sync.RWMutex
	logger       *logger.Logger
}

// NewRabbitMQEventBus 创建RabbitMQ事件总线
func NewRabbitMQEventBus(logger *logger.Logger) (EventBus, error) {
	// 从环境变量获取配置
	amqpURL := os.Getenv("RABBITMQ_URL")
	if amqpURL == "" {
		amqpURL = "amqp://guest:guest@localhost:5672/"
	}

	exchangeName := os.Getenv("RABBITMQ_EXCHANGE")
	if exchangeName == "" {
		exchangeName = "application_events"
	}

	queueName := os.Getenv("RABBITMQ_QUEUE")
	if queueName == "" {
		queueName = "application_events_queue"
	}

	// 建立连接
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// 创建通道
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	// 声明交换机
	err = ch.ExchangeDeclare(
		exchangeName, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %v", err)
	}

	// 声明队列
	queue, err := ch.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %v", err)
	}

	eventBus := &RabbitMQEventBus{
		connection:   conn,
		channel:      ch,
		exchangeName: exchangeName,
		queueName:    queue.Name,
		consumerTag:  "application-event-consumer",
		handlers:     make(map[string][]model.EventHandler),
		logger:       logger,
	}

	// 启动消费者
	go eventBus.startConsumer()

	return eventBus, nil
}

// Subscribe 订阅事件
func (rb *RabbitMQEventBus) Subscribe(eventType string, handler model.EventHandler) error {
	rb.handlerMutex.Lock()
	defer rb.handlerMutex.Unlock()

	rb.handlers[eventType] = append(rb.handlers[eventType], handler)

	// 绑定队列到交换机
	routingKey := eventType
	err := rb.channel.QueueBind(
		rb.queueName,    // queue name
		routingKey,      // routing key
		rb.exchangeName, // exchange
		false,
		nil,
	)

	if err != nil {
		return fmt.Errorf("failed to bind queue: %v", err)
	}

	return nil
}

// Publish 发布事件
func (rb *RabbitMQEventBus) Publish(ctx context.Context, event model.ApplicationEvent) error {
	// 序列化事件
	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %v", err)
	}

	// 创建消息
	msg := amqp.Publishing{
		ContentType: "application/json",
		Body:        eventData,
		Timestamp:   time.Now(),
		MessageId:   event.EventID(),
		Type:        event.EventType(),
	}

	// 发布消息
	routingKey := event.EventType()
	err = rb.channel.Publish(
		rb.exchangeName, // exchange
		routingKey,      // routing key
		false,           // mandatory
		false,           // immediate
		msg,
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	return nil
}

// Unsubscribe 取消订阅指定事件类型
func (rb *RabbitMQEventBus) Unsubscribe(eventType string, handler model.EventHandler) error {
	rb.handlerMutex.Lock()
	defer rb.handlerMutex.Unlock()

	handlers, exists := rb.handlers[eventType]
	if !exists {
		return nil
	}

	// 查找并移除对应的处理器
	for i, h := range handlers {
		if h == handler {
			rb.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}

	// 如果该事件类型没有处理器了，可以考虑解绑队列（可选）
	return nil
}

// startConsumer 启动消费者
func (rb *RabbitMQEventBus) startConsumer() {
	// 消费消息
	msgs, err := rb.channel.Consume(
		rb.queueName,   // queue
		rb.consumerTag, // consumer
		false,          // auto-ack
		false,          // exclusive
		false,          // no-local
		false,          // no-wait
		nil,            // args
	)

	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return
	}

	// 处理消息
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			// 反序列化事件
			var event map[string]interface{}
			if err := json.Unmarshal(d.Body, &event); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				d.Nack(false, false) // 拒绝消息
				continue
			}

			// 获取事件类型
			eventType := d.Type
			if eventType == "" {
				// 如果Type字段为空，尝试从事件数据中获取
				if et, ok := event["eventType"].(string); ok {
					eventType = et
				}
			}

			// 处理事件
			if err := rb.handleEvent(eventType, event); err != nil {
				log.Printf("Failed to handle event %s: %v", eventType, err)
				d.Nack(false, true) // 重新入队
				continue
			}

			// 确认消息
			d.Ack(false)
		}
	}()

	log.Printf("RabbitMQ consumer started, waiting for messages...")
	<-forever
}

// handleEvent 处理事件
func (rb *RabbitMQEventBus) handleEvent(eventType string, eventData map[string]interface{}) error {
	rb.handlerMutex.RLock()
	defer rb.handlerMutex.RUnlock()

	handlers, exists := rb.handlers[eventType]
	if !exists || len(handlers) == 0 {
		// 没有处理器，但不认为是错误
		return nil
	}

	// 创建事件对象
	event := &RabbitMQApplicationEvent{
		Data: eventData,
		Type: eventType,
	}

	// 并发处理所有处理器
	var wg sync.WaitGroup
	errChan := make(chan error, len(handlers))

	for _, handler := range handlers {
		wg.Add(1)
		go func(h model.EventHandler) {
			defer wg.Done()
			if err := h.Handle(event); err != nil {
				errChan <- err
			}
		}(handler)
	}

	wg.Wait()
	close(errChan)

	// 处理错误
	for err := range errChan {
		log.Printf("Error in event handler: %v", err)
	}

	return nil
}

// Close 关闭连接
func (rb *RabbitMQEventBus) Close() error {
	if rb.channel != nil {
		if err := rb.channel.Close(); err != nil {
			log.Printf("Error closing channel: %v", err)
		}
	}

	if rb.connection != nil {
		if err := rb.connection.Close(); err != nil {
			return fmt.Errorf("error closing connection: %v", err)
		}
	}

	return nil
}

// RabbitMQApplicationEvent RabbitMQ应用事件实现
type RabbitMQApplicationEvent struct {
	Data map[string]interface{}
	Type string
}

// EventID 获取事件ID
func (e *RabbitMQApplicationEvent) EventID() string {
	if id, ok := e.Data["id"].(string); ok {
		return id
	}
	if id, ok := e.Data["eventID"].(string); ok {
		return id
	}
	return ""
}

// EventType 获取事件类型
func (e *RabbitMQApplicationEvent) EventType() string {
	return e.Type
}

// Timestamp 获取事件时间戳
func (e *RabbitMQApplicationEvent) Timestamp() time.Time {
	if timestampStr, ok := e.Data["timestamp"].(string); ok {
		if t, err := time.Parse(time.RFC3339, timestampStr); err == nil {
			return t
		}
	}
	return time.Now()
}

// Payload 获取事件载荷
func (e *RabbitMQApplicationEvent) Payload() interface{} {
	return e.Data
}
