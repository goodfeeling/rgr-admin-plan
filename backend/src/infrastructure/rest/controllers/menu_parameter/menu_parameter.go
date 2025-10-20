package menuParameter

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMenuParameter "github.com/gbrayhan/microservices-go/src/domain/sys/menu_parameter"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewMenuParameterRequest struct {
	Type          string `json:"type"  binding:"required"`
	Key           string `json:"key"  binding:"required"`
	Value         string `json:"value"  binding:"required"`
	SysBaseMenuID int64  `json:"sys_base_menu_id"  binding:"required"`
}

type ResponseMenuParameter struct {
	ID            int    `json:"id"`
	Type          string `json:"type"`
	Key           string `json:"key"`
	Value         string `json:"value"`
	SysBaseMenuID int64  `json:"sys_base_menu_id"`
}
type IMenuParameterController interface {
	NewMenuParameter(ctx *gin.Context)
	GetAllMenuParameters(ctx *gin.Context)
	GetMenuParametersByID(ctx *gin.Context)
	UpdateMenuParameter(ctx *gin.Context)
	DeleteMenuParameter(ctx *gin.Context)
}
type MenuParameterController struct {
	menuParameterService domainMenuParameter.IMenuParameterService
	Logger               *logger.Logger
}

func NewMenuParameterController(menuParameterService domainMenuParameter.IMenuParameterService, loggerInstance *logger.Logger) IMenuParameterController {
	return &MenuParameterController{menuParameterService: menuParameterService, Logger: loggerInstance}
}

// CreateMenuParameter
// @Summary create menuParameter
// @Description create menuParameter
// @Tags menuParameter create
// @Accept json
// @Produce json
// @Param book body NewMenuParameterRequest true  "JSON Data"
// @Success 200 {object} controllers.CommonResponseBuilder
// @Router /v1/menuParameter [post]
func (c *MenuParameterController) NewMenuParameter(ctx *gin.Context) {
	c.Logger.Info("Creating new menuParameter")
	var request NewMenuParameterRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new menuParameter", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	menuParameterModel, err := c.menuParameterService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating menuParameter", zap.Error(err), zap.String("Key", request.Key))
		_ = ctx.Error(err)
		return
	}
	menuParameterResponse := controllers.NewCommonResponseBuilder[*ResponseMenuParameter]().
		Data(domainToResponseMapper(menuParameterModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("MenuParameter created successfully", zap.String("Key", request.Key), zap.Int("id", int(menuParameterModel.ID)))
	ctx.JSON(http.StatusOK, menuParameterResponse)
}

// GetAllMenuParameters
// @Summary get all menuParameters by
// @Description get  all menuParameters by where
// @Tags menuParameters
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainMenuParameter.MenuParameter]
// @Router /v1/menuParameter [get]
func (c *MenuParameterController) GetAllMenuParameters(ctx *gin.Context) {

	c.Logger.Info("Getting all menuParameters")
	c.Logger.Info("Getting all menuBtns")
	menuBaseID, err := strconv.Atoi(ctx.Query("menu_id"))
	if err != nil {
		c.Logger.Error("Invalid menuBtn ID parameter", zap.Error(err), zap.String("id", ctx.Query("menu_id")))
		appError := domainErrors.NewAppError(errors.New("menuBtn id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	menuParameters, err := c.menuParameterService.GetAll(int64(menuBaseID))
	if err != nil {
		c.Logger.Error("Error getting all menuParameters", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all menuParameters", zap.Int("count", len(*menuParameters)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]domainMenuParameter.MenuParameter]{
		Data: menuParameters,
	})
}

// GetMenuParametersByID
// @Summary get menuParameters
// @Description get menuParameters by id
// @Tags menuParameters
// @Accept json
// @Produce json
// @Success 200 {object} ResponseMenuParameter
// @Router /v1/menuParameter/{id} [get]
func (c *MenuParameterController) GetMenuParametersByID(ctx *gin.Context) {
	menuParameterID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menuParameter ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("menuParameter id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting menuParameter by ID", zap.Int("id", menuParameterID))
	menuParameter, err := c.menuParameterService.GetByID(menuParameterID)
	if err != nil {
		c.Logger.Error("Error getting menuParameter by ID", zap.Error(err), zap.Int("id", menuParameterID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved menuParameter by ID", zap.Int("id", menuParameterID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(menuParameter))
}

// UpdateMenuParameter
// @Summary update menuParameter
// @Description update menuParameter
// @Tags menuParameter
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseMenuParameter]
// @Router /v1/menuParameter [put]
func (c *MenuParameterController) UpdateMenuParameter(ctx *gin.Context) {
	menuParameterID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menuParameter ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating menuParameter", zap.Int("id", menuParameterID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for menuParameter update", zap.Error(err), zap.Int("id", menuParameterID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for menuParameter update", zap.Error(err), zap.Int("id", menuParameterID))
		_ = ctx.Error(err)
		return
	}
	menuParameterUpdated, err := c.menuParameterService.Update(menuParameterID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating menuParameter", zap.Error(err), zap.Int("id", menuParameterID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseMenuParameter]().
		Data(domainToResponseMapper(menuParameterUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("MenuParameter updated successfully", zap.Int("id", menuParameterID))
	ctx.JSON(http.StatusOK, response)
}

// DeleteMenuParameter
// @Summary delete menuParameter
// @Description delete menuParameter by id
// @Tags menuParameter
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/menuParameter/{id} [delete]
func (c *MenuParameterController) DeleteMenuParameter(ctx *gin.Context) {
	menuParameterID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menuParameter ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting menuParameter", zap.Int("id", menuParameterID))
	err = c.menuParameterService.Delete(menuParameterID)
	if err != nil {
		c.Logger.Error("Error deleting menuParameter", zap.Error(err), zap.Int("id", menuParameterID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("MenuParameter deleted successfully", zap.Int("id", menuParameterID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    menuParameterID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// Mappers
func domainToResponseMapper(domainMenuParameter *domainMenuParameter.MenuParameter) *ResponseMenuParameter {

	return &ResponseMenuParameter{
		ID:            domainMenuParameter.ID,
		Type:          domainMenuParameter.Type,
		Key:           domainMenuParameter.Key,
		Value:         domainMenuParameter.Value,
		SysBaseMenuID: domainMenuParameter.SysBaseMenuID,
	}
}
func toUsecaseMapper(req *NewMenuParameterRequest) *domainMenuParameter.MenuParameter {
	return &domainMenuParameter.MenuParameter{
		Type:          req.Type,
		Key:           req.Key,
		Value:         req.Value,
		SysBaseMenuID: req.SysBaseMenuID,
	}
}
