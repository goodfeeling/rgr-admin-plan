package di

import (
	scheduledTaskUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/scheduled_task"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/scheduled_task"
	scheduledTaskController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/scheduled_task"
)

type ScheduledTaskModule struct {
	Controller scheduledTaskController.IScheduledTaskController
	UseCase    scheduledTaskUseCase.IScheduledTaskService
	Repository scheduled_task.IScheduledTaskRepository
}

func setupScheduledTaskModule(appContext *ApplicationContext) error {

	// Initialize use cases
	service := scheduledTaskUseCase.NewScheduledTaskUseCase(
		appContext.Repositories.ScheduledTaskRepository,
		appContext.Logger, appContext.TaskScheduler)

	// Initialize controllers
	controller := scheduledTaskController.NewScheduledTaskController(
		service, appContext.Logger)

	appContext.ScheduledTaskModule = ScheduledTaskModule{
		Controller: controller,
		UseCase:    service,
		Repository: appContext.Repositories.ScheduledTaskRepository,
	}
	return nil
}
