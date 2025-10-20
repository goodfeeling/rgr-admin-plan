package scheduled_task

import (
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	scheduledTaskDomain "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
	scheduleTaskConstants "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task/constants"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/lib/scheduler"
	scheduledTaskRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/scheduled_task"
	"go.uber.org/zap"
)

type IScheduledTaskService interface {
	GetAll() (*[]scheduledTaskDomain.ScheduledTask, error)
	Create(apiDomain *scheduledTaskDomain.ScheduledTask) (*scheduledTaskDomain.ScheduledTask, error)
	GetByID(id int) (*scheduledTaskDomain.ScheduledTask, error)
	Update(id int, apiMap map[string]interface{}) (*scheduledTaskDomain.ScheduledTask, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[scheduledTaskDomain.ScheduledTask], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	EnableTask(id int) error
	DisableTask(id int) error
	ReloadTasks() error
}

type ScheduledTaskUseCase struct {
	scheduledTaskRepository scheduledTaskRepo.IScheduledTaskRepository
	Logger                  *logger.Logger
	scheduler               *scheduler.TaskScheduler
}

func NewScheduledTaskUseCase(
	scheduledTaskRepository scheduledTaskRepo.IScheduledTaskRepository,
	loggerInstance *logger.Logger, scheduler *scheduler.TaskScheduler,
) IScheduledTaskService {
	return &ScheduledTaskUseCase{
		scheduledTaskRepository: scheduledTaskRepository,
		Logger:                  loggerInstance,
		scheduler:               scheduler,
	}
}

func (s *ScheduledTaskUseCase) GetAll() (*[]scheduledTaskDomain.ScheduledTask, error) {
	s.Logger.Info("Getting all tasks")
	return s.scheduledTaskRepository.GetAll()
}

func (s *ScheduledTaskUseCase) GetByID(id int) (*scheduledTaskDomain.ScheduledTask, error) {
	s.Logger.Info("Getting task by ID", zap.Int("id", id))
	return s.scheduledTaskRepository.GetByID(id)
}

func (s *ScheduledTaskUseCase) Create(newData *scheduledTaskDomain.ScheduledTask) (*scheduledTaskDomain.ScheduledTask, error) {
	s.Logger.Info("Creating new task", zap.String("TaskName", newData.TaskName))
	return s.scheduledTaskRepository.Create(newData)
}

func (s *ScheduledTaskUseCase) Delete(ids []int) error {
	s.Logger.Info("Deleting task", zap.String("ids", fmt.Sprintf("%v", ids)))
	return s.scheduledTaskRepository.Delete(ids)
}

func (s *ScheduledTaskUseCase) Update(id int, userMap map[string]interface{}) (*scheduledTaskDomain.ScheduledTask, error) {
	s.Logger.Info("Updating task", zap.Int("id", id))
	return s.scheduledTaskRepository.Update(id, userMap)
}

func (s *ScheduledTaskUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[scheduledTaskDomain.ScheduledTask], error) {
	s.Logger.Info("Searching tasks with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.scheduledTaskRepository.SearchPaginated(filters)
}

func (s *ScheduledTaskUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching tasks by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.scheduledTaskRepository.SearchByProperty(property, searchText)
}

// DisableTask implements IScheduledTaskService.
func (s *ScheduledTaskUseCase) DisableTask(taskID int) error {
	updateData := map[string]interface{}{
		"status": scheduleTaskConstants.TaskStatusDisabled,
	}
	time.Sleep(time.Millisecond * 3000)
	_, err := s.scheduledTaskRepository.Update(taskID, updateData)
	if err != nil {
		return err
	}

	return s.scheduler.StopTask(taskID)
}

// EnableTask implements IScheduledTaskService.
func (s *ScheduledTaskUseCase) EnableTask(taskID int) error {
	updateData := map[string]interface{}{
		"status": scheduleTaskConstants.TaskStatusEnabled,
	}
	_, err := s.scheduledTaskRepository.Update(taskID, updateData)
	if err != nil {
		return err
	}

	return s.scheduler.StartTask(taskID)
}

// ReloadTasks implements IScheduledTaskService.
func (s *ScheduledTaskUseCase) ReloadTasks() error {
	return s.scheduler.ReloadTasks()
}
