// #file:/root/myproject/microapp/src/infrastructure/di/modules/auth_module.go
package di

import (
	authUseCase "github.com/gbrayhan/microservices-go/src/application/services/auth"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/jwt_blacklist"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	authController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/auth"
	wsHandler "github.com/gbrayhan/microservices-go/src/infrastructure/ws/handler/user_status"
)

type AuthModule struct {
	Controller             authController.IAuthController
	UseCase                authUseCase.IAuthUseCase
	UserRepository         user.UserRepositoryInterface
	RoleRepository         role.ISysRolesRepository
	JwtBlacklistRepository jwt_blacklist.JwtBlacklistRepository
	WsHandler              *wsHandler.UserStatusHandler
}

func setupAuthModule(appContext *ApplicationContext) error {

	// Initialize websocket handler
	wsHandler := wsHandler.NewUserStatusHandler(appContext.SessionManager, appContext.Logger)

	// Initialize use cases
	authUC := authUseCase.NewAuthUseCase(
		appContext.Repositories.UserRepository,
		appContext.Repositories.RoleRepository,
		appContext.JWTService,
		appContext.Logger,
		appContext.Repositories.JwtBlacklistRepository,
		appContext.RedisClient,
		appContext.SessionManager,
		appContext.CaptchaHandler)

	// Initialize controllers
	authController := authController.NewAuthController(authUC, appContext.Logger)

	appContext.AuthModule = AuthModule{
		Controller:             authController,
		UseCase:                authUC,
		UserRepository:         appContext.Repositories.UserRepository,
		JwtBlacklistRepository: appContext.Repositories.JwtBlacklistRepository,
		RoleRepository:         appContext.Repositories.RoleRepository,
		WsHandler:              wsHandler,
	}
	return nil
}
