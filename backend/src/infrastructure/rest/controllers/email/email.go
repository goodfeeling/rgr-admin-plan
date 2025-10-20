package email

import (
	"net/http"

	domainEmail "github.com/gbrayhan/microservices-go/src/domain/email"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IEmailController interface {
	SendForgetPasswordEmail(ctx *gin.Context)
}
type EmailController struct {
	emailService domainEmail.IEmailService
	Logger       *logger.Logger
}

// Structures
type NewForgetPasswordRequest struct {
	Email string `json:"email" binding:"required"`
}

func NewEmailController(
	emailService domainEmail.IEmailService,
	loggerInstance *logger.Logger,
) IEmailController {
	return &EmailController{emailService: emailService, Logger: loggerInstance}
}

// SendForgetPasswordEmail
// @Summary password email
// @Description forget password
// @Tags send password
// @Accept json
// @Produce json
// @Param book body models.User  true  "JSON Data"
// @Success 200 {array} models.User
// @Router /v1/email/forget-password [post]
func (e *EmailController) SendForgetPasswordEmail(ctx *gin.Context) {
	e.Logger.Info("send forget password email")
	var request NewForgetPasswordRequest
	if err := controllers.BindJSON(ctx, &request); err != nil {
		e.Logger.Error("Error binding JSON for email", zap.Error(err))
		appError := domainErrors.NewAppError(err, domainErrors.ValidationError)
		_ = ctx.Error(appError)
		return
	}
	err := e.emailService.SendForgetPasswordEmail(request.Email)
	if err != nil {
		e.Logger.Error("Error sending email", zap.Error(err), zap.String("Email", request.Email))
		appError := domainErrors.NewAppError(err, domainErrors.NotFound)
		_ = ctx.Error(appError)
		return
	}
	apiResponse := controllers.NewCommonResponseBuilder[bool]().
		Data(true).
		Message("success").
		Status(0).
		Build()

	e.Logger.Info("send email successfully", zap.String("Email", request.Email))
	ctx.JSON(http.StatusOK, apiResponse)
}
