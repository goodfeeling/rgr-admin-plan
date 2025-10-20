package di

import (
	configUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/config"
	configController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/config"
)

type ConfigModule struct {
	Controller configController.IConfigController
	UseCase    configUseCase.ISysConfigService
}

func setupConfigModule(appContext *ApplicationContext) error {
	// Initialize use cases
	configUC := configUseCase.NewSysConfigUseCase(appContext.Logger)
	// Initialize controllers
	configController := configController.NewConfigController(configUC, appContext.Logger)

	appContext.ConfigModule = ConfigModule{
		Controller: configController,
		UseCase:    configUC,
	}
	return nil
}
