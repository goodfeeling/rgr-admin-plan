package di

import (
	menuParameterUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/menu_parameter"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_parameter"
	menuParameterController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu_parameter"
)

type MenuParameterModule struct {
	Controller menuParameterController.IMenuParameterController
	UseCase    menuParameterUseCase.IMenuParameterService
	Repository base_menu_parameter.MenuParameterRepositoryInterface
}

func setupMenuParameterModule(appContext *ApplicationContext) error {
	// Initialize use cases
	menuParameterUC := menuParameterUseCase.NewMenuParameterUseCase(
		appContext.Repositories.MenuParameterRepository,
		appContext.Logger)

	// Initialize controllers
	menuParameterController := menuParameterController.NewMenuParameterController(menuParameterUC, appContext.Logger)
	appContext.MenuParameterModule = MenuParameterModule{
		Controller: menuParameterController,
		UseCase:    menuParameterUC,
		Repository: appContext.Repositories.MenuParameterRepository,
	}
	return nil

}
