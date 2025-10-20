package file

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainFile "github.com/gbrayhan/microservices-go/src/domain/sys/files"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	apiRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewFileRequest struct {
	ID             int    `json:"id"`
	FileName       string `json:"file_name" binding:"required"`
	FilePath       string `json:"file_path" binding:"required"`
	StorageEngine  string `json:"storage_engine" binding:"required"`
	FileOriginName string `json:"file_origin_name" binding:"required"`
}

// Structures
type DeleteBatchFileRequest struct {
	IDS []int64 `json:"ids"`
}

type SynchronizeResponse struct {
	Count int `json:"count"`
}

type ResponseFile struct {
	ID             int64             `json:"id"`
	FileName       string            `json:"file_name"`
	FileMD5        string            `json:"file_md5"`
	FilePath       string            `json:"file_path"`
	FileUrl        string            `json:"file_url"`
	StorageEngine  string            `json:"storage_engine"`
	FileOriginName string            `json:"file_origin_name"`
	CreatedAt      domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt      domain.CustomTime `json:"updated_at,omitempty"`
}
type IFileController interface {
	NewFile(ctx *gin.Context)
	GetAllFiles(ctx *gin.Context)
	GetFilesByID(ctx *gin.Context)
	UpdateFile(ctx *gin.Context)
	DeleteFile(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
	DeleteFiles(ctx *gin.Context)
}
type FileController struct {
	apiService domainFile.ISysFilesService
	Logger     *logger.Logger
	Router     *gin.Engine
}
type RouterSetter interface {
	SetRouter(router *gin.Engine)
}

func (c *FileController) SetRouter(router *gin.Engine) {
	c.Router = router
}

func NewFileController(apiService domainFile.ISysFilesService, loggerInstance *logger.Logger) IFileController {
	return &FileController{apiService: apiService, Logger: loggerInstance}
}

// CreateFile
// @Summary create api
// @Description create api
// @Tags api create
// @Accept json
// @Produce json
// @Param book body NewFileRequest true  "JSON Data"
// @Success 200 {object} controllers.CommonResponseBuilder
// @Router /v1/api [post]
func (c *FileController) NewFile(ctx *gin.Context) {
	c.Logger.Info("Creating new api")
	var request NewFileRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new api", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	apiModel, err := c.apiService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating api", zap.Error(err), zap.String("path", request.FileName))
		_ = ctx.Error(err)
		return
	}
	apiResponse := controllers.NewCommonResponseBuilder[*ResponseFile]().
		Data(domainToResponseMapper(apiModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("File created successfully", zap.String("path", request.FileName), zap.Int("id", int(apiModel.ID)))
	ctx.JSON(http.StatusOK, apiResponse)
}

// GetAllFiles
// @Summary get all apis by
// @Description get  all apis by where
// @Tags apis
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainFile.File]
// @Router /v1/api [get]
func (c *FileController) GetAllFiles(ctx *gin.Context) {
	c.Logger.Info("Getting all apis")
	apis, err := c.apiService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all apis", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all apis", zap.Int("count", len(*apis)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]domainFile.SysFiles]{
		Data: apis,
	})
}

// GetFilesByID
// @Summary get apis
// @Description get apis by id
// @Tags apis
// @Accept json
// @Produce json
// @Success 200 {object} ResponseFile
// @Router /v1/api/{id} [get]
func (c *FileController) GetFilesByID(ctx *gin.Context) {
	apiID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid api ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("api id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting api by ID", zap.Int("id", apiID))
	api, err := c.apiService.GetByID(apiID)
	if err != nil {
		c.Logger.Error("Error getting api by ID", zap.Error(err), zap.Int("id", apiID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved api by ID", zap.Int("id", apiID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(api))
}

// UpdateFile
// @Summary update api
// @Description update api
// @Tags api
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseFile]
// @Router /v1/api [put]
func (c *FileController) UpdateFile(ctx *gin.Context) {
	apiID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid api ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating api", zap.Int("id", apiID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for api update", zap.Error(err), zap.Int("id", apiID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for api update", zap.Error(err), zap.Int("id", apiID))
		_ = ctx.Error(err)
		return
	}
	apiUpdated, err := c.apiService.Update(apiID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating api", zap.Error(err), zap.Int("id", apiID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseFile]().
		Data(domainToResponseMapper(apiUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("File updated successfully", zap.Int("id", apiID))
	ctx.JSON(http.StatusOK, response)
}

// DeleteFile
// @Summary delete api
// @Description delete api by id
// @Tags api
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/api/{id} [delete]
func (c *FileController) DeleteFile(ctx *gin.Context) {
	apiID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid api ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting api", zap.Int("id", apiID))
	err = c.apiService.Delete([]int64{int64(apiID)})
	if err != nil {
		c.Logger.Error("Error deleting api", zap.Error(err), zap.Int("id", apiID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("File deleted successfully", zap.Int("id", apiID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    apiID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// SearchFilePageList
// @Summary search apis
// @Description search apis by query
// @Tags search apis
// @Accept json
// @Produce json
// @Success 200 {object} domain.PageList[[]ResponseFile]
// @Router /v1/api/search [get]
func (c *FileController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching apis with pagination")

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
	for field := range apiRepo.ColumnsSysFilesMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range apiRepo.ColumnsSysFilesMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}

	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range apiRepo.ColumnsSysFilesMapping {
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

	result, err := c.apiService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching apis", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseFile]
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

	c.Logger.Info("Successfully searched apis",
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
// @Router /v1/api/search-property [get]
func (c *FileController) SearchByProperty(ctx *gin.Context) {
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
		"apiName":   true,
		"email":     true,
		"firstName": true,
		"lastName":  true,
		"status":    true,
	}
	if !allowed[property] {
		c.Logger.Error("Invalid property for search", zap.String("property", property))
		appError := domainErrors.NewAppError(errors.New("invalid property"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	coincidences, err := c.apiService.SearchByProperty(property, searchText)
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

// BatchDeleteFile
// @Summary delete operations
// @Description delete operations by id
// @Tags batch delete
// @Accept json
// @Produce json
// @Param book body DeleteBatchFileRequest true  "JSON Data"
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/api/delete-batch [post]
func (c *FileController) DeleteFiles(ctx *gin.Context) {
	c.Logger.Info("Creating new dictionary")
	var request DeleteBatchFileRequest
	var err error
	if err = controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new dictionary", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting operation", zap.String("ids", fmt.Sprintf("%v", request.IDS)))
	err = c.apiService.Delete(request.IDS)
	if err != nil {
		c.Logger.Error("Error deleting operation", zap.Error(err), zap.String("ids", fmt.Sprintf("%v", request.IDS)))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("File deleted successfully", zap.String("ids", fmt.Sprintf("%v", request.IDS)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[[]int64]{
		Data:    request.IDS,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// Mappers
func domainToResponseMapper(domainFile *domainFile.SysFiles) *ResponseFile {
	return &ResponseFile{
		ID:             domainFile.ID,
		FileName:       domainFile.FileName,
		FileMD5:        domainFile.FileMD5,
		FilePath:       domainFile.FilePath,
		FileUrl:        domainFile.FileUrl,
		FileOriginName: domainFile.FileOriginName,
		StorageEngine:  domainFile.StorageEngine,
		CreatedAt:      domain.CustomTime{Time: domainFile.CreatedAt},
		UpdatedAt:      domain.CustomTime{Time: domainFile.UpdatedAt},
	}
}

func arrayDomainToResponseMapper(apis *[]domainFile.SysFiles) *[]*ResponseFile {
	res := make([]*ResponseFile, len(*apis))
	for i, u := range *apis {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewFileRequest) *domainFile.SysFiles {
	return &domainFile.SysFiles{
		FileName:       req.FileName,
		FilePath:       req.FilePath,
		FileOriginName: req.FileOriginName,
		StorageEngine:  req.StorageEngine,
	}
}
