package di

import (
	userUseCase "github.com/gbrayhan/microservices-go/src/application/services/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/job"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	userController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/user"
)

type UserModule struct {
	Controller userController.IUserController
	UseCase    userUseCase.IUserUseCase
	Repository user.UserRepositoryInterface
}

func setupUserModule(appContext *ApplicationContext) error {

	// register executor
	appContext.FunctionExecutor.RegisterFunction(job.FUNCTION_TYPE_CLEAN_UP_OLD_DATA, job.CleanOldData)

	// Initialize use cases
	userUC := userUseCase.NewUserUseCase(
		appContext.Repositories.UserRepository,
		appContext.Repositories.UserRoleRepository,
		appContext.EventBus,
		appContext.Repositories.JwtBlacklistRepository,
		appContext.Logger)

	// Initialize controllers
	userController := userController.NewUserController(userUC, appContext.Logger)

	appContext.UserModule = UserModule{
		Controller: userController,
		UseCase:    userUC,
		Repository: appContext.Repositories.UserRepository,
	}
	return nil
}
