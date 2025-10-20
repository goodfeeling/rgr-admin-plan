package menu_group

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMenuGroup "github.com/gbrayhan/microservices-go/src/domain/sys/menu_group"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	menuGroupRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_group"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewMenuGroupRequest struct {
	ID     int    `json:"id"`
	Name   string `json:"name"  binding:"required"`
	Path   string `json:"path"`
	Status int16  `json:"status"`
	Sort   int8   `json:"sort"`
}

type ResponseMenuGroup struct {
	ID        int               `json:"id"`
	Name      string            `json:"name"`
	Path      string            `json:"path"`
	Status    int16             `json:"status"`
	Sort      int8              `json:"sort"`
	CreatedAt domain.CustomTime `json:"created_at"`
	UpdatedAt domain.CustomTime `json:"updated_at"`
}
type IMenuGroupController interface {
	NewMenuGroup(ctx *gin.Context)
	GetAllMenuGroups(ctx *gin.Context)
	GetMenuGroupsByID(ctx *gin.Context)
	UpdateMenuGroup(ctx *gin.Context)
	DeleteMenuGroup(ctx *gin.Context)
	SearchPaginated(ctx *gin.Context)
	SearchByProperty(ctx *gin.Context)
}
type MenuGroupController struct {
	menuGroupService domainMenuGroup.IMenuGroupService
	Logger           *logger.Logger
}

func NewMenuGroupController(menuGroupService domainMenuGroup.IMenuGroupService, loggerInstance *logger.Logger) IMenuGroupController {
	return &MenuGroupController{menuGroupService: menuGroupService, Logger: loggerInstance}
}

// CreateMenuGroup
// @Summary create menuGroup
// @Description create menuGroup
// @Tags menuGroup create
// @Accept json
// @Produce json
// @Param book body NewMenuGroupRequest true  "JSON Data"
// @Success 200 {object} controllers.CommonResponseBuilder
// @Router /v1/menuGroup [post]
func (c *MenuGroupController) NewMenuGroup(ctx *gin.Context) {
	c.Logger.Info("Creating new menuGroup")
	var request NewMenuGroupRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new menuGroup", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	menuGroupModel, err := c.menuGroupService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating menuGroup", zap.Error(err), zap.String("Name", request.Name))
		_ = ctx.Error(err)
		return
	}
	menuGroupResponse := controllers.NewCommonResponseBuilder[*ResponseMenuGroup]().
		Data(domainToResponseMapper(menuGroupModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("MenuGroup created successfully", zap.String("Name", request.Name), zap.Int("id", int(menuGroupModel.ID)))
	ctx.JSON(http.StatusOK, menuGroupResponse)
}

// GetAllDictionaries
// @Summary get all dictionaries by
// @Description get  all dictionaries by where
// @Tags dictionaries
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainMenuGroup.MenuGroup]
// @Router /v1/menuGroup [get]
func (c *MenuGroupController) GetAllMenuGroups(ctx *gin.Context) {
	c.Logger.Info("Getting all dictionaries")
	dictionaries, err := c.menuGroupService.GetAll()
	if err != nil {
		c.Logger.Error("Error getting all dictionaries", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all dictionaries", zap.Int("count", len(*dictionaries)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]domainMenuGroup.MenuGroup]{
		Data: dictionaries,
	})
}

// GetMenuGroupsByID
// @Summary get dictionaries
// @Description get dictionaries by id
// @Tags dictionaries
// @Accept json
// @Produce json
// @Success 200 {object} ResponseMenuGroup
// @Router /v1/menuGroup/{id} [get]
func (c *MenuGroupController) GetMenuGroupsByID(ctx *gin.Context) {
	menuGroupID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menuGroup ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("menuGroup id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting menuGroup by ID", zap.Int("id", menuGroupID))
	menuGroup, err := c.menuGroupService.GetByID(menuGroupID)
	if err != nil {
		c.Logger.Error("Error getting menuGroup by ID", zap.Error(err), zap.Int("id", menuGroupID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved menuGroup by ID", zap.Int("id", menuGroupID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(menuGroup))
}

// UpdateMenuGroup
// @Summary update menuGroup
// @Description update menuGroup
// @Tags menuGroup
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseMenuGroup]
// @Router /v1/menuGroup [put]
func (c *MenuGroupController) UpdateMenuGroup(ctx *gin.Context) {
	menuGroupID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menuGroup ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating menuGroup", zap.Int("id", menuGroupID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for menuGroup update", zap.Error(err), zap.Int("id", menuGroupID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for menuGroup update", zap.Error(err), zap.Int("id", menuGroupID))
		_ = ctx.Error(err)
		return
	}
	menuGroupUpdated, err := c.menuGroupService.Update(menuGroupID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating menuGroup", zap.Error(err), zap.Int("id", menuGroupID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseMenuGroup]().
		Data(domainToResponseMapper(menuGroupUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("MenuGroup updated successfully", zap.Int("id", menuGroupID))
	ctx.JSON(http.StatusOK, response)
}

// DeleteMenuGroup
// @Summary delete menuGroup
// @Description delete menuGroup by id
// @Tags menuGroup
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/menuGroup/{id} [delete]
func (c *MenuGroupController) DeleteMenuGroup(ctx *gin.Context) {
	menuGroupID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menuGroup ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting menuGroup", zap.Int("id", menuGroupID))
	err = c.menuGroupService.Delete([]int{menuGroupID})
	if err != nil {
		c.Logger.Error("Error deleting menuGroup", zap.Error(err), zap.Int("id", menuGroupID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("MenuGroup deleted successfully", zap.Int("id", menuGroupID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    menuGroupID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// SearchMenuGroupPageList
// @Summary search dictionaries
// @Description search dictionaries by query
// @Tags search dictionaries
// @Accept json
// @Produce json
// @Success 200 {object} domain.PageList[[]ResponseMenuGroup]
// @Router /v1/menuGroup/search [get]
func (c *MenuGroupController) SearchPaginated(ctx *gin.Context) {
	c.Logger.Info("Searching dictionaries with pagination")

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
	for field := range menuGroupRepo.ColumnsMenuGroupMapping {
		if values := ctx.QueryArray(field + "_like"); len(values) > 0 {
			likeFilters[field] = values
		}
	}
	filters.LikeFilters = likeFilters

	// Parse exact matches
	matches := make(map[string][]string)
	for field := range menuGroupRepo.ColumnsMenuGroupMapping {
		if values := ctx.QueryArray(field + "_match"); len(values) > 0 {
			matches[field] = values
		}
	}
	filters.Matches = matches

	// Parse date range filters
	var dateRanges []domain.DateRangeFilter
	for field := range menuGroupRepo.ColumnsMenuGroupMapping {
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

	result, err := c.menuGroupService.SearchPaginated(filters)
	if err != nil {
		c.Logger.Error("Error searching dictionaries", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	type PageResult = domain.PageList[*[]*ResponseMenuGroup]
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

	c.Logger.Info("Successfully searched dictionaries",
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
// @Router /v1/menu-group/search-property [get]
func (c *MenuGroupController) SearchByProperty(ctx *gin.Context) {
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
		"menuGroupName": true,
		"email":         true,
		"firstName":     true,
		"lastName":      true,
		"status":        true,
	}
	if !allowed[property] {
		c.Logger.Error("Invalid property for search", zap.String("property", property))
		appError := domainErrors.NewAppError(errors.New("invalid property"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}

	coincidences, err := c.menuGroupService.SearchByProperty(property, searchText)
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

// Mappers
func domainToResponseMapper(domainMenuGroup *domainMenuGroup.MenuGroup) *ResponseMenuGroup {

	return &ResponseMenuGroup{
		ID:        domainMenuGroup.ID,
		Name:      domainMenuGroup.Name,
		Path:      domainMenuGroup.Path,
		Status:    domainMenuGroup.Status,
		Sort:      domainMenuGroup.Sort,
		CreatedAt: domain.CustomTime{Time: domainMenuGroup.CreatedAt},
		UpdatedAt: domain.CustomTime{Time: domainMenuGroup.UpdatedAt},
	}
}

func arrayDomainToResponseMapper(dictionaries *[]domainMenuGroup.MenuGroup) *[]*ResponseMenuGroup {
	res := make([]*ResponseMenuGroup, len(*dictionaries))
	for i, u := range *dictionaries {
		res[i] = domainToResponseMapper(&u)
	}
	return &res
}

func toUsecaseMapper(req *NewMenuGroupRequest) *domainMenuGroup.MenuGroup {
	return &domainMenuGroup.MenuGroup{
		Name:   req.Name,
		Path:   req.Path,
		Status: req.Status,
		Sort:   req.Sort,
	}
}
