package captcha

import (
	"net/http"

	domainCaptcha "github.com/gbrayhan/microservices-go/src/domain/captcha"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	captchaLib "github.com/gbrayhan/microservices-go/src/infrastructure/lib/captcha"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ICaptchaController interface {
	Generate(ctx *gin.Context)
	Verify(ctx *gin.Context)
}
type CaptchaController struct {
	captchaService domainCaptcha.ICaptchaService
	Logger         *logger.Logger
}

// Structures
type VerifyRequest struct {
	Id     string `json:"id" binding:"required"`
	Answer string `json:"answer" binding:"required"`
}

func NewCaptchaController(
	captchaService domainCaptcha.ICaptchaService,
	loggerInstance *logger.Logger,
) ICaptchaController {
	return &CaptchaController{captchaService: captchaService, Logger: loggerInstance}
}

// Generate implements ICaptchaController.
// @Summary captcha generate
// @Description captcha generate
// @Tags captcha
// @Accept json
// @Produce json
// @Param book body models.User  true  "JSON Data"
// @Success 200 {object} domain.CommonResponse[*captchaLib.CaptchaResponse]
// @Router /v1/captcha/generate [get]
func (c *CaptchaController) Generate(ctx *gin.Context) {

	captchaResp, err := c.captchaService.Generate(ctx.Query("captcha_id"))
	if err != nil {
		c.Logger.Error("Failed to generate captcha", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate captcha"})
		return
	}
	response := controllers.NewCommonResponseBuilder[*captchaLib.CaptchaResponse]().
		Data(captchaResp).
		Message("success").
		Status(0).
		Build()
	ctx.JSON(http.StatusOK, response)
}

// Verify
// @Summary verify captcha
// @Description verify captcha
// @Tags verify captcha
// @Accept json
// @Produce json
// @Param book body models.User  true  "JSON Data"
// @Success 200 {object} domain.CommonResponse[bool]
// @Router /v1/captcha/verify [post]
func (c *CaptchaController) Verify(ctx *gin.Context) {
	c.Logger.Info("User login request")
	var request VerifyRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		c.Logger.Error("Error binding JSON for Verify", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	isValid, err := c.captchaService.Verify(request.Id, request.Answer)
	if err != nil {
		c.Logger.Error("Failed to verify captcha", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate verify"})
		return
	}
	response := controllers.NewCommonResponseBuilder[bool]().
		Data(isValid).
		Message("success").
		Status(0).
		Build()
	ctx.JSON(http.StatusOK, response)
}
