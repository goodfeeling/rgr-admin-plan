package di

import (
	menuBtnUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/menu_btn"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_btn"
	menuBtnController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu_btn"
)

type MenuBtnModule struct {
	Controller menuBtnController.IMenuBtnController
	UseCase    menuBtnUseCase.IMenuBtnService
	Repository base_menu_btn.MenuBtnRepositoryInterface
}

func setupMenuBtnModule(appContext *ApplicationContext) error {
	// Initialize repositories
	menuBtnRepo := base_menu_btn.NewMenuBtnRepository(appContext.DB, appContext.Logger)

	// Initialize use cases
	menuBtnUC := menuBtnUseCase.NewMenuBtnUseCase(menuBtnRepo, appContext.Logger)

	// Initialize controllers
	menuBtnController := menuBtnController.NewMenuBtnController(menuBtnUC, appContext.Logger)

	appContext.MenuBtnModule = MenuBtnModule{
		Controller: menuBtnController,
		UseCase:    menuBtnUC,
		Repository: menuBtnRepo,
	}
	return nil
}
