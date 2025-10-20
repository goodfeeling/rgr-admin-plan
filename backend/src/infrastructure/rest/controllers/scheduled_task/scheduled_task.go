package scheduled_task

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainScheduledTask "github.com/gbrayhan/microservices-go/src/domain/sys/scheduled_task"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	ScheduledTaskRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/scheduled_task"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/datatypes"
)

// Structures
type DeleteBatchScheduledTaskRequest struct {
	IDS []int `json:"ids"`
}

// Structures
type NewScheduledTaskRequest struct {
	ID              int            `json:"id"`
	TaskName        string         `json:"task_name"  binding:"required"`
	TaskDescription string         `json:"task_description"  binding:"required"`
	CronExpression  string         `json:"cron_expression"  binding:"required"`
	TaskParams      datatypes.JSON `json:"task_params"`
	TaskType        string         `json:"task_type"  binding:"required"`
	ExecType        string         `json:"exec_type"  binding:"required"`
	Status          int            `json:"status"  binding:"required"`
}

type ResponseScheduledTask struct {
	ID              int               `json:"id"`
	TaskName        string            `json:"task_name"`
	TaskDescription string            `json:"task_description"`
	CronExpression  string            `json:"cron_expression"`
	TaskParams      datatypes.JSON    `json:"task_params"`
	Status          int               `json:"status"`
	TaskType        string            `json:"task_type"`
	ExecType        string            `json:"exec_type"`
	CreatedAt       domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt       domain.CustomTime `json:"updated_at,omitempty"`
	LastExecuteTime domain.CustomTime `json:"last_execute_time"`
	NextExecuteTime domain.CustomTime `json:"next_execute_time"`
}
type IScheduledTaskController interface {
	NewScheduledTask(ctx *gin.Context)
	GetAllScheduledTasks(ctx *gin.Context)
	GetScheduledTaskByID(ctx *gin.Context)
	UpdateScheduledTask(ctx *gin.Context)
	DeleteScheduledTask(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
	DeleteScheduledTasks(ctx *gin.Context)

	// task manager
	EnableTaskById(ctx *gin.Context)
	DisableTaskById(ctx *gin.Context)
	ReloadAllTasks(ctx *gin.Context)
}
type ScheduledTasController struct {
	scheduledTaskService domainScheduledTask.IScheduledTaskService
	Logger               *logger.Logger
}

func NewScheduledTaskController(
	scheduledTaskService domainScheduledTask.IScheduledTaskService,
	loggerInstance *logger.Logger,
) IScheduledTaskController {
	return &ScheduledTasController{
		scheduledTaskService: scheduledTaskService,
		Logger:               loggerInstance,
	}
}

// CreateScheduledTask
// @Summary create ScheduledTask
// @Description create ScheduledTask
// @Tags ScheduledTask create
// @Accept json
// @Produce json
// @Param book body NewScheduledTaskRequest true  "JSON Data"
// @Success 200 {object} controllers.CommonResponseBuilder
// @Router /v1/scheduled_task [post]
func (c *ScheduledTasController) NewScheduledTask(ctx *gin.Context) {
	c.Logger.Info("Creating new ScheduledTask")
	var request NewScheduledTaskRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new ScheduledTask", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	ScheduledTaskModel, err := c.scheduledTaskService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating ScheduledTask", zap.Error(err), zap.String("TaskName", request.TaskName))
		_ = ctx.Error(err)
		return
	}
	ScheduledTaskResponse := controllers.NewCommonResponseBuilder[*ResponseScheduledTask]().
		Data(domainToResponseMapper(ScheduledTaskModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("ScheduledTask created successfully", zap.String("TaskName", request.TaskName), zap.Int("id", int(ScheduledTaskModel.ID)))
	ctx.JSON(http.StatusOK, ScheduledTaskResponse)
}

// GetAllScheduledTasks
// @Summary get all scheduled_task by
// @Description get  all scheduled_task by where
// @Tags scheduled_task
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainScheduledTask.ScheduledTask]
// @Router /v1/scheduled_task [get]
func (c *ScheduledTasController) GetAllScheduledTasks(ctx *gin.Context) {
	c.Logger.Info("Getting all scheduled_task")
	scheduled_task, err := c.scheduledTaskService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all scheduled_task", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all scheduled_task", zap.Int("count", len(*scheduled_task)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]domainScheduledTask.ScheduledTask]{
		Data: scheduled_task,
	})
}

// GetScheduledTaskByID
// @Summary get scheduled_task
// @Description get scheduled_task by id
// @Tags scheduled_task
// @Accept json
// @Produce json
// @Success 200 {object} ResponseScheduledTask
// @Router /v1/scheduled_task/{id} [get]
func (c *ScheduledTasController) GetScheduledTaskByID(ctx *gin.Context) {
	ScheduledTaskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid ScheduledTask ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("ScheduledTask id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting ScheduledTask by ID", zap.Int("id", ScheduledTaskID))
	ScheduledTask, err := c.scheduledTaskService.GetByID(ScheduledTaskID)
	if err != nil {
		c.Logger.Error("Error getting ScheduledTask by ID", zap.Error(err), zap.Int("id", ScheduledTaskID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved ScheduledTask by ID", zap.Int("id", ScheduledTaskID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(ScheduledTask))
}

// UpdateScheduledTask
// @Summary update ScheduledTask
// @Description update ScheduledTask
// @Tags ScheduledTask
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseScheduledTask]
// @Router /v1/scheduled_task [put]
func (c *ScheduledTasController) UpdateScheduledTask(ctx *gin.Context) {
	ScheduledTaskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid ScheduledTask ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating ScheduledTask", zap.Int("id", ScheduledTaskID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for ScheduledTask update", zap.Error(err), zap.Int("id", ScheduledTaskID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for ScheduledTask update", zap.Error(err), zap.Int("id", ScheduledTaskID))
		_ = ctx.Error(err)
		return
	}
	ScheduledTaskUpdated, err := c.scheduledTaskService.Update(ScheduledTaskID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating ScheduledTask", zap.Error(err), zap.Int("id", ScheduledTaskID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseScheduledTask]().
		Data(domainToResponseMapper(ScheduledTaskUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("ScheduledTask updated successfully", zap.Int("id", ScheduledTaskID))
	ctx.JSON(http.StatusOK, response)
}

// DeleteScheduledTask
// @Summary delete ScheduledTask
// @Description delete ScheduledTask by id
// @Tags ScheduledTask
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/scheduled_task/{id} [delete]
func (c *ScheduledTasController) DeleteScheduledTask(ctx *gin.Context) {
	ScheduledTaskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid ScheduledTask ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting ScheduledTask", zap.Int("id", ScheduledTaskID))
	err = c.scheduledTaskService.Delete([]int{ScheduledTaskID})
	if err != nil {
		c.Logger.Error("Error deleting ScheduledTask", zap.Error(err), zap.Int("id", ScheduledTaskID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("ScheduledTask deleted successfully", zap.Int("id", ScheduledTaskID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    ScheduledTaskID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// SearchScheduledTaskPageList
// @Summary search scheduled_task
// @Description search scheduled_task by query
// @Tags search scheduled_task
// @Accept json
// @Produce json
// @Success 200 {object} domain.PageList[[]ResponseScheduledTask]
// @Router /v1/scheduled_task/search [get]
func (c *ScheduledTasController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching scheduled_task with pagination")

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
	for field := range ScheduledTaskRepo.ColumnsScheduledTaskMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range ScheduledTaskRepo.ColumnsScheduledTaskMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range ScheduledTaskRepo.ColumnsScheduledTaskMapping {
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
		c.Logger.Error("Error searching scheduled_task", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseScheduledTask]
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

	c.Logger.Info("Successfully searched scheduled_task",
		zap.Int64("total", result.Total),
		zap.Int("page", result.Page))
	ctx.JSON(http.StatusOK, response)
}

// SearchByProperty
// @Summary  search by property
// @Description search by property
// @Tags search
// @Accept json
// @Produce json
// @Success 200 {array} []string
// @Router /v1/scheduled_task/search-property [get]
func (c *ScheduledTasController) SearchByProperty(ctx *gin.Context) {
	property := ctx.Query("property")
	searchText := ctx.Query("searchText")

	if property == "" || searchText == "" {
		c.Logger.Error("Missing property or searchText parameter")
		appError := domainErrors.NewAppError(errors.New("missing property or searchText parameter"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	// Validate property
	allowed := map[string]bool{
		"ScheduledTaskName": true,
		"email":             true,
		"firstName":         true,
		"lastName":          true,
		"status":            true,
	}
	if !allowed[property] {
		c.Logger.Error("Invalid property for search", zap.String("property", property))
		appError := domainErrors.NewAppError(errors.New("invalid property"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	coincidences, err := c.scheduledTaskService.SearchByProperty(property, searchText)
	if err != nil {
		c.Logger.Error("Error searching by property", zap.Error(err), zap.String("property", property))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Successfully searched by property",
		zap.String("property", property),
		zap.Int("results", len(*coincidences)))
	ctx.JSON(http.StatusOK, coincidences)
}

// DeleteScheduledTasks
// @Summary delete tasks
// @Description delete tasks by id
// @Tags batch delete
// @Accept json
// @Produce json
// @Param book body DeleteBatchOperationRequest true  "JSON Data"
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/operation/delete-batch [post]
func (c *ScheduledTasController) DeleteScheduledTasks(ctx *gin.Context) {
	c.Logger.Info("Creating new ScheduledTask")
	var request DeleteBatchScheduledTaskRequest
	var err error
	if err = controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new ScheduledTask", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting operation", zap.String("ids", fmt.Sprintf("%v", request.IDS)))
	err = c.scheduledTaskService.Delete(request.IDS)
	if err != nil {
		c.Logger.Error("Error deleting operation", zap.Error(err), zap.String("ids", fmt.Sprintf("%v", request.IDS)))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Operation deleted successfully", zap.String("ids", fmt.Sprintf("%v", request.IDS)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[[]int]{
		Data:    request.IDS,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// Mappers
func domainToResponseMapper(domainScheduledTask *domainScheduledTask.ScheduledTask) *ResponseScheduledTask {

	return &ResponseScheduledTask{
		ID:              domainScheduledTask.ID,
		TaskName:        domainScheduledTask.TaskName,
		TaskType:        domainScheduledTask.TaskType,
		TaskParams:      domainScheduledTask.TaskParams,
		TaskDescription: domainScheduledTask.TaskDescription,
		CronExpression:  domainScheduledTask.CronExpression,
		Status:          domainScheduledTask.Status,
		LastExecuteTime: domain.CustomTime{Time: domainScheduledTask.LastExecuteTime},
		NextExecuteTime: domain.CustomTime{Time: domainScheduledTask.NextExecuteTime},
		ExecType:        domainScheduledTask.ExecType,
		CreatedAt:       domain.CustomTime{Time: domainScheduledTask.CreatedAt},
		UpdatedAt:       domain.CustomTime{Time: domainScheduledTask.UpdatedAt},
	}
}

func arrayDomainToResponseMapper(scheduled_task *[]domainScheduledTask.ScheduledTask) *[]*ResponseScheduledTask {
	res := make([]*ResponseScheduledTask, len(*scheduled_task))
	for i, u := range *scheduled_task {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewScheduledTaskRequest) *domainScheduledTask.ScheduledTask {
	return &domainScheduledTask.ScheduledTask{
		CronExpression:  req.CronExpression,
		Status:          req.Status,
		TaskDescription: req.TaskDescription,
		TaskName:        req.TaskName,
		TaskParams:      req.TaskParams,
		TaskType:        req.TaskType,
		ExecType:        req.ExecType,
	}
}

// ReloadAllTasks implements IScheduledTaskController.
// @Summary reload all tasks
// @Description reload all tasks
// @Tags tasks
// @Accept json
// @Produce json
// @Param book body models.User  true  "JSON Data"
// @Success 200 {array} models.User
// @Router /api/v1/scheduled_task/reload-all [post]
func (c *ScheduledTasController) ReloadAllTasks(ctx *gin.Context) {
	c.Logger.Info("Reload all ScheduledTask ")
	err := c.scheduledTaskService.ReloadTasks()
	if err != nil {
		c.Logger.Error("Error reload ScheduledTask")
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully Reload ScheduledTask")
	ctx.JSON(http.StatusOK, domain.CommonResponse[bool]{
		Data:    true,
		Message: "resource reload successfully",
		Status:  0,
	})
}

// EnableTask implements IScheduledTaskController.
// @Summary start task
// @Description enable task
// @Tags task
// @Accept json
// @Produce json
// @Param book body models.User  true  "JSON Data"
// @Success 200 {array} models.User
// @Router /api/v1/scheduled_task/enable/{id} [post]
func (c *ScheduledTasController) EnableTaskById(ctx *gin.Context) {
	scheduledTaskID, err := strconv.Atoi(ctx.Param("id"))

	if err != nil {
		c.Logger.Error("Invalid ScheduledTask ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("ScheduledTask id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Starting ScheduledTask by ID", zap.Int("id", scheduledTaskID))
	err = c.scheduledTaskService.EnableTask(scheduledTaskID)
	if err != nil {
		c.Logger.Error("Error Starting ScheduledTask by ID", zap.Error(err), zap.Int("id", scheduledTaskID))
		_ = ctx.Error(err)
		return
	}

	c.Logger.Info("Successfully Starting ScheduledTask by ID", zap.Int("id", scheduledTaskID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    scheduledTaskID,
		Message: "resource start successfully",
		Status:  0,
	})
}

// StopTask implements IScheduledTaskController.
// @Summary disable task
// @Description disable task
// @Tags disable task
// @Accept json
// @Produce json
// @Param book body models.User  true  "JSON Data"
// @Success 200 {array} models.User
// @Router /api/v1/scheduled_task/disable/{id} [post]
func (c *ScheduledTasController) DisableTaskById(ctx *gin.Context) {
	scheduledTaskID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid ScheduledTask ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("ScheduledTask id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("disabling ScheduledTask by ID", zap.Int("id", scheduledTaskID))
	err = c.scheduledTaskService.DisableTask(scheduledTaskID)
	if err != nil {
		c.Logger.Error("Error disabling ScheduledTask by ID", zap.Error(err), zap.Int("id", scheduledTaskID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully disabling ScheduledTask by ID", zap.Int("id", scheduledTaskID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    scheduledTaskID,
		Message: "resource disabled successfully",
		Status:  0,
	})
}
