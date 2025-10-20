package di

import (
	apiUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/api"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/api"

	apiController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/api"
)

type ApiModule struct {
	Controller apiController.IApiController
	UseCase    apiUseCase.ISysApiService
	Repository api.ApiRepositoryInterface
}

func setupApiModule(appContext *ApplicationContext) error {
	// Initialize repositories
	apiRepository := api.NewApiRepository(appContext.DB, appContext.Logger)
	// Initialize use cases
	apiUC := apiUseCase.NewSysApiUseCase(
		apiRepository,
		appContext.Repositories.DictionaryRepository,
		appContext.Logger)
	// Initialize controllers
	apiController := apiController.NewApiController(apiUC, appContext.Logger)

	appContext.ApiModule = ApiModule{
		Controller: apiController,
		UseCase:    apiUC,
		Repository: apiRepository,
	}
	return nil
}
