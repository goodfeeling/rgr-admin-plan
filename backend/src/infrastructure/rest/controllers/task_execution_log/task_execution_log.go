package task_execution_log

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainTaskExecutionLog "github.com/gbrayhan/microservices-go/src/domain/sys/task_execution_log"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	TaskExecutionLogRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/task_execution_log"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ResponseTaskExecutionLog struct {
	ID              int               `json:"id"`
	TaskID          uint              `json:"task_id"`
	ExecuteTime     time.Time         `json:"execute_time"`
	ExecuteResult   int               `json:"execute_result"`
	ExecuteDuration *int              `json:"execute_duration"`
	ErrorMessage    string            `json:"error_message"`
	CreatedAt       domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt       domain.CustomTime `json:"updated_at,omitempty"`
}

type ITaskExecutionLogController interface {
	SearchPaginated(ctx *gin.Context)
}

type TaskExecutionLogController struct {
	scheduledTaskService domainTaskExecutionLog.ITaskExecutionLogService
	Logger               *logger.Logger
}

func NewTaskExecutionLogController(
	scheduledTaskService domainTaskExecutionLog.ITaskExecutionLogService,
	loggerInstance *logger.Logger,
) ITaskExecutionLogController {
	return &TaskExecutionLogController{
		scheduledTaskService: scheduledTaskService,
		Logger:               loggerInstance,
	}
}

// SearchTaskExecutionLogPageList
// @Summary search task_execution_log
// @Description search task_execution_log by query
// @Tags search task_execution_log
// @Accept json
// @Produce json
// @Success 200 {object} domain.PageList[[]ResponseTaskExecutionLog]
// @Router /v1/task_execution_log/search [get]
func (c *TaskExecutionLogController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching task_execution_log with pagination")

	// Parse query parameters
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))
	if pageSize < 1 {
		pageSize = 10
	}

	// Build filters
	filters := domain.DataFilters{
		Page:     page,
		PageSize: pageSize,
	}

	// Parse like filters
	likeFilters := make(map[string][]string)
	for field := range TaskExecutionLogRepo.ColumnsTaskExecutionLogMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range TaskExecutionLogRepo.ColumnsTaskExecutionLogMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range TaskExecutionLogRepo.ColumnsTaskExecutionLogMapping {
		startStr := ctx.Query(field + "_start")
		endStr := ctx.Query(field + "_end")

		if startStr != "" || endStr != "" {
			dateRange := domain.DateRangeFilter{Field: field}

			if startStr != "" {
				if startTime, err := time.Parse(time.RFC3339, startStr); err == nil {
					dateRange.Start = &startTime
				}
			}

			if endStr != "" {
				if endTime, err := time.Parse(time.RFC3339, endStr); err == nil {
					dateRange.End = &endTime
				}
			}

			dateRanges = append(dateRanges, dateRange)
		}
	}
	filters.DateRangeFilters = dateRanges

	// Parse sorting
	sortBy := ctx.QueryArray("sortBy")
	if len(sortBy) > 0 {
		filters.SortBy = sortBy
	} else {
		filters.SortBy = []string{}
	}

	sortDirection := domain.SortDirection(ctx.DefaultQuery("sortDirection", "asc"))
	if sortDirection.IsValid() {
		filters.SortDirection = sortDirection
	}

	result, err := c.scheduledTaskService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching task_execution_log", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseTaskExecutionLog]
	response := controllers.NewCommonResponseBuilder[PageResult]().
		Data(PageResult{
			List:       arrayDomainToResponseMapper(result.Data),
			Total:      result.Total,
			Page:       result.Page,
			PageSize:   result.PageSize,
			TotalPages: result.TotalPages,
			Filters:    filters,
		}).
		Message("success").
		Status(0).
		Build()

	c.Logger.Info("Successfully searched task_execution_log",
		zap.Int64("total", result.Total),
		zap.Int("page", result.Page))
	ctx.JSON(http.StatusOK, response)
}

// Mappers
func domainToResponseMapper(domainTaskExecutionLog *domainTaskExecutionLog.TaskExecutionLog) *ResponseTaskExecutionLog {
	return &ResponseTaskExecutionLog{
		ID:              domainTaskExecutionLog.ID,
		TaskID:          domainTaskExecutionLog.TaskID,
		ExecuteTime:     domainTaskExecutionLog.ExecuteTime,
		ExecuteResult:   domainTaskExecutionLog.ExecuteResult,
		ExecuteDuration: domainTaskExecutionLog.ExecuteDuration,
		ErrorMessage:    domainTaskExecutionLog.ErrorMessage,
		CreatedAt:       domain.CustomTime{Time: domainTaskExecutionLog.CreatedAt},
		UpdatedAt:       domain.CustomTime{Time: domainTaskExecutionLog.UpdatedAt},
	}
}

func arrayDomainToResponseMapper(task_execution_log *[]domainTaskExecutionLog.TaskExecutionLog) *[]*ResponseTaskExecutionLog {
	res := make([]*ResponseTaskExecutionLog, len(*task_execution_log))
	for i, u := range *task_execution_log {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}
