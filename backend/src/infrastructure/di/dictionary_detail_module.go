package di

import (
	dictionaryDetailUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/dictionary_detail"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary_detail"
	dictionaryDetailController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/dictionary_detail"
)

type DictionaryDetailModule struct {
	Controller dictionaryDetailController.IIDictionaryDetailController
	UseCase    dictionaryDetailUseCase.ISysDictionaryService
	Repository dictionary_detail.DictionaryRepositoryInterface
}

func setupDictionaryDetailModule(appContext *ApplicationContext) error {
	// Initialize repositories
	dictionaryDetailRepo := dictionary_detail.NewDictionaryRepository(appContext.DB, appContext.Logger)

	// Initialize use cases
	dictionaryDetailUC := dictionaryDetailUseCase.NewSysDictionaryUseCase(dictionaryDetailRepo, appContext.Logger)

	// Initialize controllers
	dictionaryDetailController := dictionaryDetailController.NewIDictionaryDetailController(dictionaryDetailUC, appContext.Logger)

	appContext.DictionaryDetailModule = DictionaryDetailModule{
		Controller: dictionaryDetailController,
		UseCase:    dictionaryDetailUC,
		Repository: dictionaryDetailRepo,
	}
	return nil
}
