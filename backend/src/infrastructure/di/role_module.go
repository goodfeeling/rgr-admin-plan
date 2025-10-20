package di

import (
	roleUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/casbin_rule"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_btn"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_menu"
	roleController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/role"
)

type RoleModule struct {
	Controller         roleController.IRoleController
	UseCase            roleUseCase.ISysRoleService
	Repository         role.ISysRolesRepository
	RoleMenuRepository role_menu.ISysRoleMenuRepository
	CasBinRepository   casbin_rule.ICasbinRuleRepository
	MenuRepository     base_menu.MenuRepositoryInterface
	RoleBtnRepository  role_btn.ISysRoleBtnRepository
}

func setupRoleModule(appContext *ApplicationContext) error {

	// Initialize use cases
	roleUC := roleUseCase.NewSysRoleUseCase(
		appContext.Repositories.RoleRepository,
		appContext.Repositories.RoleMenuRepository,
		appContext.Repositories.CasBinRepository,
		appContext.Repositories.MenuRepository,
		appContext.Repositories.RoleBtnRepository,
		appContext.Logger)

	// Initialize controllers
	roleController := roleController.NewRoleController(roleUC, appContext.Logger)
	appContext.RoleModule = RoleModule{
		Controller:         roleController,
		UseCase:            roleUC,
		Repository:         appContext.Repositories.RoleRepository,
		RoleMenuRepository: appContext.Repositories.RoleMenuRepository,
		RoleBtnRepository:  appContext.Repositories.RoleBtnRepository,
		CasBinRepository:   appContext.Repositories.CasBinRepository,
		MenuRepository:     appContext.Repositories.MenuRepository,
	}
	return nil
}
