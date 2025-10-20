package di

import (
	dictionaryUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/dictionary"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary"
	dictionaryController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/dictionary"
)

type DictionaryModule struct {
	Controller dictionaryController.IDictionaryController
	UseCase    dictionaryUseCase.ISysDictionaryService
	Repository dictionary.DictionaryRepositoryInterface
}

func setupDictionaryModule(appContext *ApplicationContext) error {
	// Initialize repositories
	dictionaryRepo := dictionary.NewDictionaryRepository(appContext.DB, appContext.Logger)

	// Initialize use cases
	dictionaryUC := dictionaryUseCase.NewSysDictionaryUseCase(dictionaryRepo, appContext.Logger)

	// Initialize controllers
	dictionaryController := dictionaryController.NewDictionaryController(dictionaryUC, appContext.Logger)
	appContext.DictionaryModule = DictionaryModule{
		Controller: dictionaryController,
		UseCase:    dictionaryUC,
		Repository: dictionaryRepo,
	}
	return nil

}
