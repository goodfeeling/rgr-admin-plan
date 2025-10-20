package scheduled_task

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"go.uber.org/zap"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ScheduledTask struct {
	ID              int            `gorm:"primaryKey" json:"id"`
	TaskName        string         `gorm:"size:150;not null;uniqueIndex" json:"task_name"`
	TaskDescription string         `json:"task_description"`
	CronExpression  string         `gorm:"size:50;not null" json:"cron_expression"`
	TaskType        string         `gorm:"size:100;not null" json:"task_type"`
	TaskParams      datatypes.JSON `json:"task_params"`
	ExecType        string         `json:"exec_type"`
	Status          int            `gorm:"default:1" json:"status"`
	LastExecuteTime time.Time      `json:"last_execute_time"`
	NextExecuteTime time.Time      `json:"next_execute_time"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (ScheduledTask) TableName() string {
	return "sys_scheduled_tasks"
}

var ColumnsScheduledTaskMapping = map[string]string{
	"id":       "id",
	"status":   "status",
	"taskType": "task_type",
	"taskName": "task_name",
}

type IScheduledTaskRepository interface {
	GetAll() (*[]domainScheduledTask.ScheduledTask, error)
	Create(domainScheduledTask *domainScheduledTask.ScheduledTask) (*domainScheduledTask.ScheduledTask, error)
	GetByID(id int) (*domainScheduledTask.ScheduledTask, error)
	Update(id int, taskMap map[string]interface{}) (*domainScheduledTask.ScheduledTask, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainScheduledTask.ScheduledTask], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
}
type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewScheduledTaskRepository(db *gorm.DB, loggerInstance *logger.Logger) IScheduledTaskRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll() (*[]domainScheduledTask.ScheduledTask, error) {
	var tasks []ScheduledTask
	if err := r.DB.Find(&tasks).Error; err != nil {
		r.Logger.Error("Error getting all tasks", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all tasks", zap.Int("count", len(tasks)))
	return arrayToDomainMapper(&tasks), nil
}

func (r *Repository) Create(domainScheduledTask *domainScheduledTask.ScheduledTask) (*domainScheduledTask.ScheduledTask, error) {
	r.Logger.Info("Creating new task", zap.String("TaskName", domainScheduledTask.TaskName))
	taskRepository := fromDomainMapper(domainScheduledTask)
	txDb := r.DB.Create(taskRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating task", zap.Error(err), zap.String("TaskName", domainScheduledTask.TaskName))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return nil, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return nil, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created task", zap.String("TaskName", domainScheduledTask.TaskName), zap.Int("id", int(taskRepository.ID)))
	return taskRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainScheduledTask.ScheduledTask, error) {
	var task ScheduledTask
	err := r.DB.Where("id = ?", id).First(&task).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("data not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting task by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainScheduledTask.ScheduledTask{}, err
	}
	r.Logger.Info("Successfully retrieved task by ID", zap.Int("id", id))
	return task.toDomainMapper(), nil
}

func (r *Repository) Update(id int, dataMap map[string]interface{}) (*domainScheduledTask.ScheduledTask, error) {
	var dataObj ScheduledTask
	dataObj.ID = id
	delete(dataMap, "updated_at")
	err := r.DB.Model(&dataObj).Updates(dataMap).Error
	if err != nil {
		r.Logger.Error("Error updating task", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainScheduledTask.ScheduledTask{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &domainScheduledTask.ScheduledTask{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &domainScheduledTask.ScheduledTask{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&dataObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated task", zap.Error(err), zap.Int("id", id))
		return &domainScheduledTask.ScheduledTask{}, err
	}
	r.Logger.Info("Successfully updated task", zap.Int("id", id))
	return dataObj.toDomainMapper(), nil
}

func (r *Repository) Delete(ids []int) error {
	tx := r.DB.Where("id IN ?", ids).Delete(&ScheduledTask{})

	if tx.Error != nil {
		r.Logger.Error("Error deleting task", zap.Error(tx.Error), zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("data not found for deletion", zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted task", zap.String("ids", fmt.Sprintf("%v", ids)))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainScheduledTask.ScheduledTask], error) {
	query := r.DB.Model(&ScheduledTask{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsScheduledTaskMapping[field]
					if column != "" {
						query = query.Where(column+" ILIKE ?", "%"+value+"%")
					}
				}
			}
		}
	}

	// Apply exact matches
	for field, values := range filters.Matches {
		if len(values) > 0 {
			column := ColumnsScheduledTaskMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsScheduledTaskMapping[dateFilter.Field]
		if column != "" {
			if dateFilter.Start != nil {
				query = query.Where(column+" >= ?", dateFilter.Start)
			}
			if dateFilter.End != nil {
				query = query.Where(column+" <= ?", dateFilter.End)
			}
		}
	}

	// Apply sorting
	if len(filters.SortBy) > 0 && filters.SortDirection.IsValid() {
		for _, sortField := range filters.SortBy {
			column := ColumnsScheduledTaskMapping[sortField]
			if column != "" {
				query = query.Order(column + " " + string(filters.SortDirection))
			}
		}
	}

	// Count total records
	var total int64
	clonedQuery := query
	clonedQuery.Count(&total)

	// Apply pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10
	}
	offset := (filters.Page - 1) * filters.PageSize

	var tasks []ScheduledTask
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&tasks).Error; err != nil {
		r.Logger.Error("Error searching tasks", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainScheduledTask.ScheduledTask]{
		Data:       arrayToDomainMapper(&tasks),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched tasks",
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) SearchByProperty(property string, searchText string) (*[]string, error) {
	column := ColumnsScheduledTaskMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&ScheduledTask{}).
		Distinct(column).
		Where(column+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(column, &coincidences).Error; err != nil {
		r.Logger.Error("Error searching by property", zap.Error(err), zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	r.Logger.Info("Successfully searched by property",
		zap.String("property", property),
		zap.Int("results", len(coincidences)))

	return &coincidences, nil
}

func (u *ScheduledTask) toDomainMapper() *domainScheduledTask.ScheduledTask {
	return &domainScheduledTask.ScheduledTask{
		ID:              u.ID,
		TaskName:        u.TaskName,
		TaskType:        u.TaskType,
		TaskDescription: u.TaskDescription,
		TaskParams:      u.TaskParams,
		CronExpression:  u.CronExpression,
		Status:          u.Status,
		ExecType:        u.ExecType,
		LastExecuteTime: u.LastExecuteTime,
		NextExecuteTime: u.NextExecuteTime,

		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func fromDomainMapper(u *domainScheduledTask.ScheduledTask) *ScheduledTask {
	return &ScheduledTask{
		ID:              u.ID,
		TaskName:        u.TaskName,
		TaskType:        u.TaskType,
		TaskDescription: u.TaskDescription,
		TaskParams:      u.TaskParams,
		CronExpression:  u.CronExpression,
		Status:          u.Status,
		ExecType:        u.ExecType,
		LastExecuteTime: u.LastExecuteTime,
		NextExecuteTime: u.NextExecuteTime,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func arrayToDomainMapper(tasks *[]ScheduledTask) *[]domainScheduledTask.ScheduledTask {
	tasksDomain := make([]domainScheduledTask.ScheduledTask, len(*tasks))
	for i, task := range *tasks {
		tasksDomain[i] = *task.toDomainMapper()
	}
	return &tasksDomain
}
