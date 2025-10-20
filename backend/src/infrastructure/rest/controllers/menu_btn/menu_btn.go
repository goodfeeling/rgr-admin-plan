package menuBtn

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMenuBtn "github.com/gbrayhan/microservices-go/src/domain/sys/menu_btn"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewMenuBtnRequest struct {
	Name          string `json:"name"  binding:"required"`
	Desc          string `json:"desc"  binding:"required"`
	SysBaseMenuID int64  `json:"sys_base_menu_id"  binding:"required"`
}

type ResponseMenuBtn struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"  binding:"required"`
	Desc          string            `json:"desc"  binding:"required"`
	SysBaseMenuID int64             `json:"sys_base_menu_id"  binding:"required"`
	CreatedAt     domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt     domain.CustomTime `json:"updated_at,omitempty"`
}
type IMenuBtnController interface {
	NewMenuBtn(ctx *gin.Context)
	GetAllMenuBtns(ctx *gin.Context)
	GetMenuBtnsByID(ctx *gin.Context)
	UpdateMenuBtn(ctx *gin.Context)
	DeleteMenuBtn(ctx *gin.Context)
}
type MenuBtnController struct {
	menuBtnService domainMenuBtn.IMenuBtnService
	Logger         *logger.Logger
}

func NewMenuBtnController(menuBtnService domainMenuBtn.IMenuBtnService, loggerInstance *logger.Logger) IMenuBtnController {
	return &MenuBtnController{menuBtnService: menuBtnService, Logger: loggerInstance}
}

// CreateMenuBtn
// @Summary create menuBtn
// @Description create menuBtn
// @Tags menuBtn create
// @Accept json
// @Produce json
// @Param book body NewMenuBtnRequest true  "JSON Data"
// @Success 200 {object} controllers.CommonResponseBuilder
// @Router /v1/menuBtn [post]
func (c *MenuBtnController) NewMenuBtn(ctx *gin.Context) {
	c.Logger.Info("Creating new menuBtn")
	var request NewMenuBtnRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new menuBtn", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	menuBtnModel, err := c.menuBtnService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating menuBtn", zap.Error(err), zap.String("Name", request.Name))
		_ = ctx.Error(err)
		return
	}
	menuBtnResponse := controllers.NewCommonResponseBuilder[*ResponseMenuBtn]().
		Data(domainToResponseMapper(menuBtnModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("MenuBtn created successfully", zap.String("Name", request.Name), zap.Int("id", int(menuBtnModel.ID)))
	ctx.JSON(http.StatusOK, menuBtnResponse)
}

// GetAllMenuBtns
// @Summary get all menuBtns by
// @Description get  all menuBtns by where
// @Tags menuBtns
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainMenuBtn.MenuBtn]
// @Router /v1/menu_btn [get]
func (c *MenuBtnController) GetAllMenuBtns(ctx *gin.Context) {
	c.Logger.Info("Getting all menuBtns")
	menuBaseID, err := strconv.Atoi(ctx.Query("menu_id"))
	if err != nil {
		c.Logger.Error("Invalid menuBtn ID parameter", zap.Error(err), zap.String("id", ctx.Query("menu_id")))
		appError := domainErrors.NewAppError(errors.New("menuBtn id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	menuBtns, err := c.menuBtnService.GetAll(int64(menuBaseID))
	if err != nil {
		c.Logger.Error("Error getting all menuBtns", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all menuBtns", zap.Int("count", len(*menuBtns)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*[]domainMenuBtn.MenuBtn]{
		Data: menuBtns,
	})
}

// GetMenuBtnsByID
// @Summary get menuBtns
// @Description get menuBtns by id
// @Tags menuBtns
// @Accept json
// @Produce json
// @Success 200 {object} ResponseMenuBtn
// @Router /v1/menuBtn/{id} [get]
func (c *MenuBtnController) GetMenuBtnsByID(ctx *gin.Context) {
	menuBtnID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menuBtn ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("menuBtn id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting menuBtn by ID", zap.Int("id", menuBtnID))
	menuBtn, err := c.menuBtnService.GetByID(menuBtnID)
	if err != nil {
		c.Logger.Error("Error getting menuBtn by ID", zap.Error(err), zap.Int("id", menuBtnID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved menuBtn by ID", zap.Int("id", menuBtnID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(menuBtn))
}

// UpdateMenuBtn
// @Summary update menuBtn
// @Description update menuBtn
// @Tags menuBtn
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseMenuBtn]
// @Router /v1/menuBtn [put]
func (c *MenuBtnController) UpdateMenuBtn(ctx *gin.Context) {
	menuBtnID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menuBtn ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating menuBtn", zap.Int("id", menuBtnID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for menuBtn update", zap.Error(err), zap.Int("id", menuBtnID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for menuBtn update", zap.Error(err), zap.Int("id", menuBtnID))
		_ = ctx.Error(err)
		return
	}
	menuBtnUpdated, err := c.menuBtnService.Update(menuBtnID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating menuBtn", zap.Error(err), zap.Int("id", menuBtnID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseMenuBtn]().
		Data(domainToResponseMapper(menuBtnUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("MenuBtn updated successfully", zap.Int("id", menuBtnID))
	ctx.JSON(http.StatusOK, response)
}

// DeleteMenuBtn
// @Summary delete menuBtn
// @Description delete menuBtn by id
// @Tags menuBtn
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/menuBtn/{id} [delete]
func (c *MenuBtnController) DeleteMenuBtn(ctx *gin.Context) {
	menuBtnID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid menuBtn ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting menuBtn", zap.Int("id", menuBtnID))
	err = c.menuBtnService.Delete(menuBtnID)
	if err != nil {
		c.Logger.Error("Error deleting menuBtn", zap.Error(err), zap.Int("id", menuBtnID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("MenuBtn deleted successfully", zap.Int("id", menuBtnID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    menuBtnID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// Mappers
func domainToResponseMapper(domainMenuBtn *domainMenuBtn.MenuBtn) *ResponseMenuBtn {

	return &ResponseMenuBtn{
		ID:            domainMenuBtn.ID,
		Name:          domainMenuBtn.Name,
		Desc:          domainMenuBtn.Desc,
		SysBaseMenuID: domainMenuBtn.SysBaseMenuID,
	}
}

func toUsecaseMapper(req *NewMenuBtnRequest) *domainMenuBtn.MenuBtn {
	return &domainMenuBtn.MenuBtn{
		Name:          req.Name,
		Desc:          req.Desc,
		SysBaseMenuID: req.SysBaseMenuID,
	}
}
