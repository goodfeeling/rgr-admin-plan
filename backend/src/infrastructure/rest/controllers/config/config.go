package config

import (
	"net/http"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainConfig "github.com/gbrayhan/microservices-go/src/domain/sys/config"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ResponseConfig struct {
	ID          int64             `json:"id"`
	ConfigKey   string            `json:"config_key"`
	ConfigValue string            `json:"config_value"`
	ConfigType  string            `json:"config_type"`
	Module      string            `json:"module"`
	Sort        int               `json:"sort"`
	CreatedAt   domain.CustomTime `json:"created_at,omitempty"`
	UpdatedAt   domain.CustomTime `json:"updated_at,omitempty"`
}
type IConfigController interface {
	GetAllConfigs(ctx *gin.Context)
	UpdateConfig(ctx *gin.Context)
	GetConfigByModule(ctx *gin.Context)
	GetConfigBySite(ctx *gin.Context)
}
type ConfigController struct {
	configService domainConfig.IConfigService
	Logger        *logger.Logger
}

func NewConfigController(configService domainConfig.IConfigService, loggerInstance *logger.Logger) IConfigController {
	return &ConfigController{configService: configService, Logger: loggerInstance}
}

// GetAllConfigs
// @Summary get all configs by
// @Description get  all configs by where
// @Tags configs
// @Accept json
// @Produce json
// @Success 200 {object} domain.CommonResponse[[]domainConfig.Config]
// @Router /v1/api/config [get]
func (c *ConfigController) GetAllConfigs(ctx *gin.Context) {
	c.Logger.Info("Getting all configs")
	configs, err := c.configService.GetConfig()
	if err != nil {
		c.Logger.Error("Error getting all configs", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all configs")
	ctx.JSON(http.StatusOK, domain.CommonResponse[*domainConfig.ConfigResponse]{
		Data: configs,
	})
}

// UpdateConfig
// @Summary update config
// @Description update config
// @Tags config
// @Accept json
// @Produce json
// @Param book body map[string]any  true  "JSON Data"
// @Success 200 {array} controllers.CommonResponseBuilder[ResponseConfig]
// @Router /v1/api/config/{module} [put]
func (c *ConfigController) UpdateConfig(ctx *gin.Context) {
	var requestMap map[string]any
	err := controllers.BindJSONMap(ctx, &requestMap)
	if err != nil {
		c.Logger.Error("Error binding JSON for config update", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err = updateValidation(requestMap)
	if err != nil {
		c.Logger.Error("Validation error for config update", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	err = c.configService.Update(ctx.Param("module"), requestMap)
	if err != nil {
		c.Logger.Error("Error updating config", zap.Error(err))
		_ = ctx.Error(err)
		return
	}
	response := controllers.NewCommonResponseBuilder[bool]().
		Data(true).
		Message("success").
		Status(0).
		Build()
	c.Logger.Info("Config updated successfully")
	ctx.JSON(http.StatusOK, response)
}

// GetConfigByModule implements IConfigController.
// @Summary config module
// @Description config module
// @Tags config
// @Accept json
// @Produce json
// @Param book body models.User  true  "JSON Data"
// @Success 200 {array} models.User
// @Router /v1/api/config/module/{module} [get]
func (c *ConfigController) GetConfigByModule(ctx *gin.Context) {
	c.Logger.Info("Getting all configs")
	configs, err := c.configService.GetConfigByModule(ctx.Param("module"))
	if err != nil {
		c.Logger.Error("Error getting all configs", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all configs", zap.Int("count", len(*configs)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*map[string]string]{
		Data: configs,
	})
}

// GetConfigBySiteConfig implements IConfigController.
// @Summary config module
// @Description config module
// @Tags config
// @Accept json
// @Produce json
// @Param book body models.User  true  "JSON Data"
// @Success 200 {array} models.User
// @Router /v1/api/config/system [get]
func (c *ConfigController) GetConfigBySite(ctx *gin.Context) {
	c.Logger.Info("Getting all configs")
	configs, err := c.configService.GetConfigByModule("site")
	if err != nil {
		c.Logger.Error("Error getting all configs", zap.Error(err))
		appError := domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		_ = ctx.Error(appError)
		return
	}
	c.Logger.Info("Successfully retrieved all configs", zap.Int("count", len(*configs)))
	ctx.JSON(http.StatusOK, domain.CommonResponse[*map[string]string]{
		Data: configs,
	})
}
