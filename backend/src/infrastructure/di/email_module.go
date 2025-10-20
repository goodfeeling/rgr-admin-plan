package di

import (
	eventHandler "github.com/gbrayhan/microservices-go/src/application/event/handler"
	eventModel "github.com/gbrayhan/microservices-go/src/application/event/model"
	emailUseCase "github.com/gbrayhan/microservices-go/src/application/services/email"
	emailController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/email"
)

type EmailModule struct {
	Controller emailController.IEmailController
	UseCase    emailUseCase.IEmailService
}

func setupEmailModule(appContext *ApplicationContext) error {
	// Initialize event
	appContext.EventBus.Subscribe(eventModel.ForgetPasswordEventType, eventHandler.NewEmailEventHandler())

	// Initialize use cases
	service := emailUseCase.NewEmailUseCase(
		appContext.Repositories.UserRepository,
		appContext.JWTService,
		appContext.EventBus,
		appContext.RedisClient,
		appContext.Logger)

	// Initialize controllers
	controller := emailController.NewEmailController(service, appContext.Logger)

	appContext.EmailModule = EmailModule{
		Controller: controller,
		UseCase:    service,
	}
	return nil
}
