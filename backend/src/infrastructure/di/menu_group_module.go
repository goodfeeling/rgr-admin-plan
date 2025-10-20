package di

import (
	menuGroupUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/menu_group"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_group"
	menuGroupController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu_group"
)

type MenuGroupModule struct {
	Controller menuGroupController.IMenuGroupController
	UseCase    menuGroupUseCase.ISysMenuGroupService
	Repository base_menu_group.MenuGroupRepositoryInterface
}

func setupMenuGroupModule(appContext *ApplicationContext) error {

	// Initialize use cases
	menuGroupUC := menuGroupUseCase.NewSysMenuGroupUseCase(
		appContext.Repositories.MenuGroupRepository, appContext.Logger)

	// Initialize controllers
	menuGroupController := menuGroupController.NewMenuGroupController(menuGroupUC, appContext.Logger)

	appContext.MenuGroupModule = MenuGroupModule{
		Controller: menuGroupController,
		UseCase:    menuGroupUC,
		Repository: appContext.Repositories.MenuGroupRepository,
	}
	return nil
}
