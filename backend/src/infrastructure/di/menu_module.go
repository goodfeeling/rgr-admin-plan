package di

import (
	menuUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_group"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_btn"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	menuController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/menu"
)

type MenuModule struct {
	Controller          menuController.IMenuController
	UseCase             menuUseCase.ISysMenuService
	Repository          base_menu.MenuRepositoryInterface
	RoleMenuRepository  role_menu.ISysRoleMenuRepository
	UserRepository      user.UserRepositoryInterface
	MenuGroupRepository base_menu_group.MenuGroupRepositoryInterface
	RoleBtnRepository   role_btn.ISysRoleBtnRepository
}

func setupMenuModule(appContext *ApplicationContext) error {
	// Initialize use cases
	menuUC := menuUseCase.NewSysMenuUseCase(
		appContext.Repositories.MenuRepository,
		appContext.Repositories.RoleMenuRepository,
		appContext.Repositories.UserRepository,
		appContext.Repositories.MenuGroupRepository,
		appContext.Repositories.RoleBtnRepository, appContext.Logger)

	// Initialize controllers
	menuController := menuController.NewMenuController(menuUC, appContext.Logger)
	appContext.MenuModule = MenuModule{
		Controller:          menuController,
		UseCase:             menuUC,
		Repository:          appContext.Repositories.MenuRepository,
		RoleMenuRepository:  appContext.Repositories.RoleMenuRepository,
		UserRepository:      appContext.Repositories.UserRepository,
		MenuGroupRepository: appContext.Repositories.MenuGroupRepository,
		RoleBtnRepository:   appContext.Repositories.RoleBtnRepository,
	}
	return nil
}
