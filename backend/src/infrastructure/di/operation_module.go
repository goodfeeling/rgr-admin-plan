package di

import (
	"github.com/gbrayhan/microservices-go/src/application/services/sys/operation_record"
	operationUseCase "github.com/gbrayhan/microservices-go/src/application/services/sys/operation_record"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/operation_records"
	operationController "github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/operation"
)

type OperationModule struct {
	Controller operationController.IOperationController
	UseCase    operationUseCase.ISysOperationService
	Repository operation_record.ISysOperationService
}

func setupOperationModule(appContext *ApplicationContext) error {
	// Initialize repositories
	operationRepo := operation_records.NewOperationRepository(appContext.DB, appContext.Logger)

	// Initialize use cases
	operationUC := operationUseCase.NewSysOperationUseCase(operationRepo, appContext.Logger)

	// Initialize controllers
	operationController := operationController.NewOperationController(operationUC, appContext.Logger)
	appContext.OperationModule = OperationModule{
		Controller: operationController,
		UseCase:    operationUC,
		Repository: operationRepo,
	}
	return nil

}
