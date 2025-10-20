package di

import (
	taskExecutionLogUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/task_execution_log"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/task_execution_log"
	taskExecutionLogController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/task_execution_log"
	wsHandler "github.com/gbrayhan/microservices-go/src/infrastructure/ws/handler/task_execution_log"
)

type TaskExecutionLogModule struct {
	Controller taskExecutionLogController.ITaskExecutionLogController
	UseCase    taskExecutionLogUseCase.ITaskExecutionLogService
	Repository task_execution_log.ITaskExecutionLogRepository
	WsHandler  *wsHandler.LogHandler
}

func setupTaskExecutionLogModule(appContext *ApplicationContext) error {

	// Initialize use cases
	services := taskExecutionLogUseCase.NewTaskExecutionLogUseCase(
		appContext.Repositories.TaskExecutionLogRepository,
		appContext.Logger)

	// Initialize websocket handler
	wsHandler := wsHandler.NewLogHandler(services, appContext.Logger)

	// Initialize controllers
	ctrl := taskExecutionLogController.NewTaskExecutionLogController(services, appContext.Logger)

	appContext.TaskExecutionLogModule = TaskExecutionLogModule{
		Controller: ctrl,
		UseCase:    services,
		Repository: appContext.Repositories.TaskExecutionLogRepository,
		WsHandler:  wsHandler,
	}
	return nil
}
