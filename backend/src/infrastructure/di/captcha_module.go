package di

import (
	captchaUseCase "github.com/gbrayhan/microservices-go/src/application/services/captcha"
	captchaController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/captcha"
)

type CaptchaModule struct {
	Controller captchaController.ICaptchaController
	UseCase    captchaUseCase.ICaptchaService
}

func setupCaptchaModule(appContext *ApplicationContext) error {

	// Initialize use cases
	service := captchaUseCase.NewCaptchaUseCase(appContext.Logger, appContext.CaptchaHandler)

	// Initialize controllers
	controller := captchaController.NewCaptchaController(service, appContext.Logger)

	appContext.CaptchaModule = CaptchaModule{
		Controller: controller,
		UseCase:    service,
	}
	return nil
}
