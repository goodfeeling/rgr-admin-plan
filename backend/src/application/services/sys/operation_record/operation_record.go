package operation_record

import (
	"fmt"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	operationRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/operation_records"

	"github.com/gbrayhan/microservices-go/src/domain"
	operationDomain "github.com/gbrayhan/microservices-go/src/domain/sys/operation_records"
	"go.uber.org/zap"
)

type ISysOperationService interface {
	GetAll() (*[]operationDomain.SysOperationRecord, error)
	GetByID(id int) (*operationDomain.SysOperationRecord, error)
	Create(newOperation *operationDomain.SysOperationRecord) (*operationDomain.SysOperationRecord, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*operationDomain.SysOperationRecord, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[operationDomain.SysOperationRecord], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*operationDomain.SysOperationRecord, error)
}

type SysOperationUseCase struct {
	sysOperationRepository operationRepo.OperationRepositoryInterface
	Logger                 *logger.Logger
}

func NewSysOperationUseCase(sysOperationRepository operationRepo.OperationRepositoryInterface, loggerInstance *logger.Logger) ISysOperationService {
	return &SysOperationUseCase{
		sysOperationRepository: sysOperationRepository,
		Logger:                 loggerInstance,
	}
}

func (s *SysOperationUseCase) GetAll() (*[]operationDomain.SysOperationRecord, error) {
	s.Logger.Info("Getting all roles")
	return s.sysOperationRepository.GetAll()
}

func (s *SysOperationUseCase) GetByID(id int) (*operationDomain.SysOperationRecord, error) {
	s.Logger.Info("Getting operation by ID", zap.Int("id", id))
	return s.sysOperationRepository.GetByID(id)
}

func (s *SysOperationUseCase) Create(newOperation *operationDomain.SysOperationRecord) (*operationDomain.SysOperationRecord, error) {
	s.Logger.Info("Creating new operation", zap.String("path", newOperation.Path))
	return s.sysOperationRepository.Create(newOperation)
}

func (s *SysOperationUseCase) Delete(ids []int) error {
	s.Logger.Info("Deleting operation", zap.String("ids", fmt.Sprintf("%v", ids)))
	return s.sysOperationRepository.Delete(ids)
}

func (s *SysOperationUseCase) Update(id int, userMap map[string]interface{}) (*operationDomain.SysOperationRecord, error) {
	s.Logger.Info("Updating operation", zap.Int("id", id))
	return s.sysOperationRepository.Update(id, userMap)
}

func (s *SysOperationUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[operationDomain.SysOperationRecord], error) {
	s.Logger.Info("Searching operations with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysOperationRepository.SearchPaginated(filters)
}

func (s *SysOperationUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching operation by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysOperationRepository.SearchByProperty(property, searchText)
}

func (s *SysOperationUseCase) GetOneByMap(userMap map[string]interface{}) (*operationDomain.SysOperationRecord, error) {
	return s.sysOperationRepository.GetOneByMap(userMap)
}
