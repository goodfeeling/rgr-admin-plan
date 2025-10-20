package di

import (
	filesUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/files"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
	uploadController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/upload"
)

type UploadModule struct {
	Controller uploadController.IUploadController
	UseCase    filesUseCase.ISysFilesService
	Repository files.ISysFilesRepository
}

func setupUploadModule(appContext *ApplicationContext) error {
	// Initialize use cases
	service := filesUseCase.NewSysFilesUseCase(appContext.Repositories.FileRepository, appContext.Logger)

	// Initialize controllers
	controller := uploadController.NewAuthController(service, appContext.Logger, appContext.RedisClient)
	appContext.UploadModule = UploadModule{
		Controller: controller,
		UseCase:    service,
		Repository: appContext.Repositories.FileRepository,
	}
	return nil
}
