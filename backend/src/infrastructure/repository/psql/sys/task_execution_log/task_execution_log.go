package task_execution_log

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainTaskExecution "github.com/gbrayhan/microservices-go/src/domain/sys/task_execution_log"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TaskExecutionLog struct {
	ID              int       `gorm:"primaryKey" json:"id"`
	TaskID          uint      `gorm:"not null;index" json:"task_id"`
	ExecuteTime     time.Time `gorm:"not null;index" json:"execute_time"`
	ExecuteResult   int       `gorm:"not null" json:"execute_result"` // 1-成功, 0-失败
	ExecuteDuration *int      `json:"execute_duration"`               // 执行耗时(毫秒)
	ErrorMessage    string    `json:"error_message"`
	CreatedAt       time.Time `json:"created_at"`
}

// TableName 指定表名
func (TaskExecutionLog) TableName() string {
	return "sys_task_execution_logs"
}

type ITaskExecutionLogRepository interface {
	GetAll() (*[]domainTaskExecution.TaskExecutionLog, error)
	Create(logDomain *domainTaskExecution.TaskExecutionLog) (*domainTaskExecution.TaskExecutionLog, error)
	GetByID(id int) (*domainTaskExecution.TaskExecutionLog, error)
	Update(id int, apiMap map[string]interface{}) (*domainTaskExecution.TaskExecutionLog, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainTaskExecution.TaskExecutionLog], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetByTaskID(taskID uint, limit int) (*[]domainTaskExecution.TaskExecutionLog, error)
}

var ColumnsTaskExecutionLogMapping = map[string]string{
	"taskId": "task_id",
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewTaskExecutionLogRepository(db *gorm.DB, loggerInstance *logger.Logger) ITaskExecutionLogRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll() (*[]domainTaskExecution.TaskExecutionLog, error) {
	var tasks []TaskExecutionLog
	if err := r.DB.Find(&tasks).Error; err != nil {
		r.Logger.Error("Error getting all tasks", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all tasks", zap.Int("count", len(tasks)))
	return arrayToDomainMapper(&tasks), nil
}

func (r *Repository) Create(domainTaskExecutionLog *domainTaskExecution.TaskExecutionLog) (*domainTaskExecution.TaskExecutionLog, error) {
	r.Logger.Info("Creating new api", zap.Uint("TaskID", domainTaskExecutionLog.TaskID))
	logRepository := fromDomainMapper(domainTaskExecutionLog)
	txDb := r.DB.Create(logRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating api", zap.Error(err), zap.Uint("TaskID", domainTaskExecutionLog.TaskID))
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
	r.Logger.Info("Successfully created api", zap.Uint("TaskID", domainTaskExecutionLog.TaskID), zap.Int("id", int(logRepository.ID)))
	return logRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainTaskExecution.TaskExecutionLog, error) {
	var api TaskExecutionLog
	err := r.DB.Where("id = ?", id).First(&api).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Api not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting api by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainTaskExecution.TaskExecutionLog{}, err
	}
	r.Logger.Info("Successfully retrieved api by ID", zap.Int("id", id))
	return api.toDomainMapper(), nil
}

func (r *Repository) Update(id int, dataMap map[string]interface{}) (*domainTaskExecution.TaskExecutionLog, error) {
	var dataObj TaskExecutionLog
	dataObj.ID = id
	delete(dataMap, "updated_at")
	err := r.DB.Model(&dataObj).Updates(dataMap).Error
	if err != nil {
		r.Logger.Error("Error updating api", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainTaskExecution.TaskExecutionLog{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &domainTaskExecution.TaskExecutionLog{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &domainTaskExecution.TaskExecutionLog{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&dataObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated api", zap.Error(err), zap.Int("id", id))
		return &domainTaskExecution.TaskExecutionLog{}, err
	}
	r.Logger.Info("Successfully updated api", zap.Int("id", id))
	return dataObj.toDomainMapper(), nil
}

func (r *Repository) Delete(ids []int) error {
	tx := r.DB.Where("id IN ?", ids).Delete(&TaskExecutionLog{})

	if tx.Error != nil {
		r.Logger.Error("Error deleting api", zap.Error(tx.Error), zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Api not found for deletion", zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted api", zap.String("ids", fmt.Sprintf("%v", ids)))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainTaskExecution.TaskExecutionLog], error) {
	query := r.DB.Model(&TaskExecutionLog{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsTaskExecutionLogMapping[field]
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
			column := ColumnsTaskExecutionLogMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsTaskExecutionLogMapping[dateFilter.Field]
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
			column := ColumnsTaskExecutionLogMapping[sortField]
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

	var tasks []TaskExecutionLog
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&tasks).Error; err != nil {
		r.Logger.Error("Error searching tasks", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainTaskExecution.TaskExecutionLog]{
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
	column := ColumnsTaskExecutionLogMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&TaskExecutionLog{}).
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

// GetByTaskID implements ITaskExecutionLogRepository.
func (r *Repository) GetByTaskID(taskID uint, limit int) (*[]domainTaskExecution.TaskExecutionLog, error) {
	var tasks []TaskExecutionLog
	if err := r.DB.Where("task_id = ?", taskID).Order("ID desc").Limit(limit).Find(&tasks).Error; err != nil {
		r.Logger.Error("Error retrieving tasks by taskID", zap.Error(err), zap.Uint("taskID", taskID))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return arrayToDomainMapper(&tasks), nil
}

func (u *TaskExecutionLog) toDomainMapper() *domainTaskExecution.TaskExecutionLog {
	return &domainTaskExecution.TaskExecutionLog{
		ID:              u.ID,
		TaskID:          u.TaskID,
		ExecuteResult:   u.ExecuteResult,
		ExecuteTime:     u.ExecuteTime,
		ErrorMessage:    u.ErrorMessage,
		ExecuteDuration: u.ExecuteDuration,
		CreatedAt:       u.CreatedAt,
	}
}

func fromDomainMapper(u *domainTaskExecution.TaskExecutionLog) *TaskExecutionLog {
	return &TaskExecutionLog{
		ID:              u.ID,
		TaskID:          u.TaskID,
		ExecuteResult:   u.ExecuteResult,
		ExecuteTime:     u.ExecuteTime,
		ErrorMessage:    u.ErrorMessage,
		ExecuteDuration: u.ExecuteDuration,
		CreatedAt:       u.CreatedAt,
	}
}

func arrayToDomainMapper(tasks *[]TaskExecutionLog) *[]domainTaskExecution.TaskExecutionLog {
	tasksDomain := make([]domainTaskExecution.TaskExecutionLog, len(*tasks))
	for i, api := range *tasks {
		tasksDomain[i] = *api.toDomainMapper()
	}
	return &tasksDomain
}
