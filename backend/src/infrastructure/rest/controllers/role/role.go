package role

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainRole "github.com/gbrayhan/microservices-go/src/domain/sys/role"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Structures
type NewRoleRequest struct {
	ID          int64  `json:"id"`
	Name        string `json:"name" binding:"required"`
	ParentID    int64  `json:"parent_id"`
	Order       int64  `json:"order"`
	Label       string `json:"label"`
	Status      int16  `json:"status"`
	Description string `json:"description"`
}

type ResponseRole struct {
	ID            int64             `json:"id"`
	Name          string            `json:"name"`
	ParentID      int64             `json:"parent_id"`
	Order         int64             `json:"order"`
	Label         string            `json:"label"`
	Status        int16             `json:"status"`
	Description   string            `json:"description"`
	DefaultRouter string            `json:"default_router"`
	CreatedAt     domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt     domain.CustomTime `json:"updated_at,omitempty"`
}

type MenuRoleResponse struct {
	RoleMenus map[int][]int     `json:"role_menus"`
	RoleBtns  map[int64][]int64 `json:"role_btns"`
	RoleApis  []string          `json:"role_apis"`
}
type IRoleController interface {
	NewRole(ctx *gin.Context)
	GetAllRoles(ctx *gin.Context)
	GetRolesByID(ctx *gin.Context)
	UpdateRole(ctx *gin.Context)
	DeleteRole(ctx *gin.Context)
	GetTreeRoles(ctx *gin.Context)
	GetRoleSetting(ctx *gin.Context)
	UpdateRoleMenuIds(ctx *gin.Context)
	BindApiRule(ctx *gin.Context)
	BindRoleMenuBtns(ctx *gin.Context)
}
type RoleController struct {
	roleService domainRole.IRoleService
	Logger      *logger.Logger
}

func NewRoleController(roleService domainRole.IRoleService, loggerInstance *logger.Logger) IRoleController {
	return &RoleController{roleService: roleService, Logger: loggerInstance}
}

// CreateRole
// @Summary create role
// @Description create role
// @Tags role create
// @Accept json
// @Produce json
// @Param book body NewRoleRequest true  "JSON Data"
// @Success 200 {object} controllers.CommonResponseBuilder
// @Router /v1/role [post]
func (c *RoleController) NewRole(ctx *gin.Context) {
	c.Logger.Info("Creating new role")
	var request NewRoleRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for new role", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	roleModel, err := c.roleService.Create(toUsecaseMapper(&request))
	if err != nil {
		c.Logger.Error("Error creating role", zap.Error(err), zap.String("Name", request.Name))
		_ = ctx.Error(err)
		return
	}
	roleResponse := controllers.NewCommonResponseBuilder[*ResponseRole]().
		Data(domainToResponseMapper(roleModel)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role created successfully", zap.String("Name", request.Name), zap.Int("id", int(roleModel.ID)))
	ctx.JSON(http.StatusOK, roleResponse)
}

// GetAllRoles
// @Summary get all roles by
// @Description get  all roles by where
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainRole.RoleTree]
// @Router /v1/role [get]
func (c *RoleController) GetAllRoles(ctx *gin.Context) {
	c.Logger.Info("Getting all roles")
	status, err := strconv.Atoi(ctx.Query("status"))
	if err != nil {
		status = 0
	}
	roles, err := c.roleService.GetAll(status)
	if err != nil {
		c.Logger.Error("Error getting all roles", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	response := controllers.NewCommonResponseBuilder[[]*domainRole.RoleTree]().
		Data(roles).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Successfully retrieved all roles", zap.Int("count", len(roles)))
	ctx.JSON(http.StatusOK, response)
}

// GetRolesByID
// @Summary get roles
// @Description get roles by id
// @Tags roles
// @Accept json
// @Produce json
// @Success 200 {object} ResponseRole
// @Router /v1/role/{id} [get]
func (c *RoleController) GetRolesByID(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("role id is invalid"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Getting role by ID", zap.Int("id", roleID))
	role, err := c.roleService.GetByID(roleID)
	if err != nil {
		c.Logger.Error("Error getting role by ID", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Successfully retrieved role by ID", zap.Int("id", roleID))
	ctx.JSON(http.StatusOK, domainToResponseMapper(role))
}

// UpdateRole
// @Summary update role
// @Description update role
// @Tags role
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseRole]
// @Router /v1/role [put]
func (c *RoleController) UpdateRole(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating role", zap.Int("id", roleID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for role update", zap.Error(err), zap.Int("id", roleID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for role update", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	roleUpdated, err := c.roleService.Update(roleID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating role", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[*ResponseRole]().
		Data(domainToResponseMapper(roleUpdated)).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role updated successfully", zap.Int("id", roleID))
	ctx.JSON(http.StatusOK, response)
}

// DeleteRole
// @Summary delete role
// @Description delete role by id
// @Tags role
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[int]
// @Router /v1/role/{id} [delete]
func (c *RoleController) DeleteRole(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter for deletion", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Deleting role", zap.Int("id", roleID))
	err = c.roleService.Delete(roleID)
	if err != nil {
		c.Logger.Error("Error deleting role", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	c.Logger.Info("Role deleted successfully", zap.Int("id", roleID))
	ctx.JSON(http.StatusOK, domain.CommonResponse[int]{
		Data:    roleID,
		Message: "resource deleted successfully",
		Status:  0,
	})
}

// GetTreeRoles
// @Summary get tree roles
// @Description get tree roles
// @Tags tree roles
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[domainRole.RoleNode]
// @Router /v1/role/tree [get]
func (c *RoleController) GetTreeRoles(ctx *gin.Context) {
	c.Logger.Info("Getting all roles tree")
	status, err := strconv.Atoi(ctx.Query("status"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter for get role menu ids", zap.Error(err), zap.String("id", ctx.Query("status")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	roles, err := c.roleService.GetTreeRoles(status)
	if err != nil {
		c.Logger.Error("Error getting all roles tree", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all roles tree", zap.Int("count", len(roles.Children)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*domainRole.RoleNode]{
		Data: roles,
	})
}

// GetRoleSetting
func (c *RoleController) GetRoleSetting(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter for get role menu ids", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	roleMenus, roleBtns, err := c.roleService.GetRoleMenuIds(int64(roleID))
	if err != nil {
		c.Logger.Error("Error getting role menu ids", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}

	if roleMenus == nil {
		roleMenus = make(map[int][]int)
	}
	rules, err := c.roleService.GetApiRuleList(roleID)
	if err != nil {
		c.Logger.Error("Error getting all apis", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	if rules == nil {
		rules = []string{}
	}
	response := controllers.NewCommonResponseBuilder[MenuRoleResponse]().
		Data(MenuRoleResponse{
			RoleMenus: roleMenus,
			RoleBtns:  roleBtns,
			RoleApis:  rules,
		}).
		Message("success").
		Status(0).
		Build()
	ctx.JSON(http.StatusOK, response)
}

func (c *RoleController) UpdateRoleMenuIds(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter for update", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Updating role", zap.Int("id", roleID))
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for role update", zap.Error(err), zap.Int("id", roleID))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = c.roleService.UpdateRoleMenuIds(roleID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating role", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[bool]().
		Data(true).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role updated successfully", zap.Int("id", roleID))
	ctx.JSON(http.StatusOK, response)
}

// BindApiRule
// @Summary bind api rule
// @Description casbin rule bind
// @Tags bind casbin
// @Accept json
// @Produce json
// @Param book body models.User  true  "JSON Data"
// @Success 200 {array} models.User
// @Router /v1/role/{id}/api [post]
func (c *RoleController) BindApiRule(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter for get role menu ids", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for casbin rule update", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = c.roleService.BindApiRule(roleID, requestMap)
	if err != nil {
		c.Logger.Error("Error updating casbin rule ", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[bool]().
		Data(true).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role updated successfully", zap.Int("id", roleID))
	ctx.JSON(http.StatusOK, response)
}

// Bind RoleMenuBtns
// @Summary bind role menu buttons
// @Description bind role menu buttons
// @Tags buttons menu role
// @Accept json
// @Produce json
// @Success 200 {array} models.User
// @Router /v1/role/{id}/menu-btns [post]
func (c *RoleController) BindRoleMenuBtns(ctx *gin.Context) {
	roleID, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		c.Logger.Error("Invalid role ID parameter for get role menu ids", zap.Error(err), zap.String("id", ctx.Param("id")))
		appError := domainErrors.NewAppError(errors.New("param id is necessary"), domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	var requestMap map[string]any
	err = controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for casbin rule update", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = c.roleService.BindRoleMenuBtns(int64(roleID), requestMap)
	if err != nil {
		c.Logger.Error("Error updating casbin rule ", zap.Error(err), zap.Int("id", roleID))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[bool]().
		Data(true).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Role updated successfully", zap.Int("id", roleID))
	ctx.JSON(http.StatusOK, response)
}

// Mappers
func domainToResponseMapper(domainRole *domainRole.Role) *ResponseRole {
	return &ResponseRole{
		ID:            domainRole.ID,
		Name:          domainRole.Name,
		ParentID:      domainRole.ParentID,
		Order:         domainRole.Order,
		Label:         domainRole.Label,
		Description:   domainRole.Description,
		Status:        domainRole.Status,
		DefaultRouter: domainRole.DefaultRouter,
		CreatedAt:     domainRole.CreatedAt,
		UpdatedAt:     domainRole.UpdatedAt,
	}
}

func toUsecaseMapper(req *NewRoleRequest) *domainRole.Role {
	return &domainRole.Role{
		Name:        req.Name,
		ParentID:    req.ParentID,
		Description: req.Description,
		Order:       req.Order,
		Label:       req.Label,
		Status:      req.Status,
	}
}
