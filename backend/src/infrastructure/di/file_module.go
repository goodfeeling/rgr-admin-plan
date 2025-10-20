package di

import (
	filesUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/files"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
	fileController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/file"
)

type FileModule struct {
	Controller fileController.IFileController
	UseCase    filesUseCase.ISysFilesService
	Repository files.ISysFilesRepository
}

func setupFileModule(appContext *ApplicationContext) error {
	// Initialize use cases
	filesUC := filesUseCase.NewSysFilesUseCase(
		appContext.Repositories.FileRepository,
		appContext.Logger)

	// Initialize controllers
	fileController := fileController.NewFileController(filesUC, appContext.Logger)

	appContext.FileModule = FileModule{
		Controller: fileController,
		UseCase:    filesUC,
		Repository: appContext.Repositories.FileRepository,
	}
	return nil
}
