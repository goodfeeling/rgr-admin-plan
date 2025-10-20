package operation

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainOperation "github.com/gbrayhan/microservices-go/src/domain/sys/operation_records"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	operationRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/operation_records"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type DeleteBatchOperationRequest struct {
	IDS []int `json:"ids"`
}

type ResponseOperation struct {
	ID           int               `json:"id"`
	IP           string            `json:"ip"`
	Path         string            `json:"path"`
	Method       string            `json:"method"`
	Status       int64             `json:"status"`
	Latency      int64             `json:"latency"`
	Agent        string            `json:"agent"`
	ErrorMessage string            `json:"error_message"`
	Body         string            `json:"body"`
	Resp         string            `json:"resp"`
	CreatedAt    domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt    domain.CustomTime `json:"updated_at,omitempty"`
}
type IOperationController interface {
	GetAllOperations(ctx *gin.Context)
	GetOperationsByID(ctx *gin.Context)
	DeleteOperation(ctx *gin.Context)
	DeleteOperations(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
}
type OperationController struct {
	operationService domainOperation.ISysOperationRecordService
	Logger           *logger.Logger
}

func NewOperationController(operationService domainOperation.ISysOperationRecordService, loggerInstance *logger.Logger) IOperationController {
	return &OperationController{operationService: operationService, Logger: loggerInstance}
}

// GetAllOperations
// @Summary get all operations by
// @Description get  all operations by where
// @Tags operations
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainOperation.Operation]
// @Router /v1/operation [get]
func (c *OperationController) GetAllOperations(ctx *gin.Context) {
	c.Logger.Info("Getting all operations")
	operations, err := c.operationService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all operations", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all operations", zap.Int("count", len(*operations)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]domainOperation.SysOperationRecord]{
		Data: operations,
	})
}

// GetOperationsByID
// @Summary get operations
// @Description get operations by id
// @Tags operations
// @Accept json
// @Produce json
// @Success 200 {object} ResponseOperation
// @Router /v1/operation/{id} [get]
func (c *OperationController) GetOperationsByID(ctx *gin.Context) {
	operationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid operation ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("operation id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting operation by ID", zap.Int("id", operationID))
	operation, err := c.operationService.GetByID(operationID)
	if err != nil {
		c.Logger.Error("Error getting operation by ID", zap.Error(err), zap.Int("id", operationID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved operation by ID", zap.Int("id", operationID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(operation))
}

// DeleteOperation
// @Summary delete operation
// @Description delete operation by id
// @Tags operation
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/operation/{id} [delete]
func (c *OperationController) DeleteOperation(ctx *gin.Context) {
	operationID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid operation ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting operation", zap.Int("id", operationID))
	err = c.operationService.Delete([]int{operationID})
	if err != nil {
		c.Logger.Error("Error deleting operation", zap.Error(err), zap.Int("id", operationID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Operation deleted successfully", zap.Int("id", operationID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    operationID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// BatchDeleteOperation
// @Summary delete operations
// @Description delete operations by id
// @Tags batch delete
// @Accept json
// @Produce json
// @Param book body DeleteBatchOperationRequest true  "JSON Data"
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/operation/delete-batch [post]
func (c *OperationController) DeleteOperations(ctx *gin.Context) {
	c.Logger.Info("Creating new dictionary")
	var request DeleteBatchOperationRequest
	var err error
	if err = controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new dictionary", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting operation", zap.String("ids", fmt.Sprintf("%v", request.IDS)))
	err = c.operationService.Delete(request.IDS)
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

// SearchOperationPageList
// @Summary search operations
// @Description search operations by query
// @Tags search operations
// @Accept json
// @Produce json
// @Success 200 {object} domain.PageList[[]ResponseOperation]
// @Router /v1/operation/search [get]
func (c *OperationController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching operations with pagination")

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
	for field := range operationRepo.ColumnsOperationMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range operationRepo.ColumnsOperationMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range operationRepo.ColumnsOperationMapping {
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

	result, err := c.operationService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching operations", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseOperation]
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

	c.Logger.Info("Successfully searched operations",
		zap.Int64("total", result.Total),
		zap.Int("page", result.Page))
	ctx.JSON(http.StatusOK, response)
}

// Mappers
func domainToResponseMapper(domainOperation *domainOperation.SysOperationRecord) *ResponseOperation {

	return &ResponseOperation{
		ID:           domainOperation.ID,
		IP:           domainOperation.IP,
		Path:         domainOperation.Path,
		Method:       domainOperation.Method,
		Status:       domainOperation.Status,
		Latency:      domainOperation.Latency,
		Agent:        domainOperation.Agent,
		ErrorMessage: domainOperation.ErrorMessage,
		Body:         domainOperation.Body,
		Resp:         domainOperation.Resp,
		CreatedAt:    domainOperation.CreatedAt,
		UpdatedAt:    domainOperation.UpdatedAt,
	}
}

func arrayDomainToResponseMapper(operations *[]domainOperation.SysOperationRecord) *[]*ResponseOperation {
	res := make([]*ResponseOperation, len(*operations))
	for i, u := range *operations {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}
