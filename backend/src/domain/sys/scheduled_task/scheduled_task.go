package scheduled_task

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	"gorm.io/datatypes"
)

type ScheduledTask struct {
	ID              int            `json:"id"`
	TaskName        string         `json:"task_name"`
	TaskDescription string         `json:"task_description"`
	CronExpression  string         `json:"cron_expression"`
	TaskType        string         `json:"task_type"`
	TaskParams      datatypes.JSON `json:"task_params"`
	Status          int            `json:"status"`
	ExecType        string         `json:"exec_type"`
	LastExecuteTime time.Time      `json:"last_execute_time"`
	NextExecuteTime time.Time      `json:"next_execute_time"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

type IScheduledTaskService interface {
	GetAll() (*[]ScheduledTask, error)
	Create(apiDomain *ScheduledTask) (*ScheduledTask, error)
	GetByID(id int) (*ScheduledTask, error)
	Update(id int, apiMap map[string]interface{}) (*ScheduledTask, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[ScheduledTask], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	EnableTask(id int) error
	DisableTask(id int) error
	ReloadTasks() error
}
