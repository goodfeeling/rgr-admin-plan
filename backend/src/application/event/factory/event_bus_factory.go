package factory

import (
	"os"

	"github.com/gbrayhan/microservices-go/src/application/event/bus"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"go.uber.org/zap"
)

// CreateEventBus 根据环境创建事件总线
func CreateEventBus(logger *logger.Logger) bus.EventBus {
	eventBusType := os.Getenv("SERVER_EVENT_BUS")

	switch eventBusType {
	case "watermill":
		// return CreateWatermillEventBus()
		return nil
	case "rabbitmq":
		rabbitBus, err := bus.NewRabbitMQEventBus(logger)
		if err != nil {
			logger.Error("Failed to create RabbitMQ event bus: %v, falling back to in-memory", zap.Error(err))
			return bus.NewInMemoryEventBus(logger)
		}
		return rabbitBus
	default:
		return bus.NewInMemoryEventBus(logger)
	}
}
