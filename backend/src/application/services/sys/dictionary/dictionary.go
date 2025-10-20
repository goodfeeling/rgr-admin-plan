package dictionary

import (
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	dictionaryRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary"

	"github.com/gbrayhan/microservices-go/src/domain"
	dictionaryDomain "github.com/gbrayhan/microservices-go/src/domain/sys/dictionary"
	"go.uber.org/zap"
)

type ISysDictionaryService interface {
	GetAll() (*[]dictionaryDomain.Dictionary, error)
	GetByID(id int) (*dictionaryDomain.Dictionary, error)
	Create(newDictionary *dictionaryDomain.Dictionary) (*dictionaryDomain.Dictionary, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*dictionaryDomain.Dictionary, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[dictionaryDomain.Dictionary], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*dictionaryDomain.Dictionary, error)
	GetByType(typeText string) (*dictionaryDomain.Dictionary, error)
}

type SysDictionaryUseCase struct {
	sysDictionaryRepository dictionaryRepo.DictionaryRepositoryInterface
	Logger                  *logger.Logger
}

func NewSysDictionaryUseCase(
	sysDictionaryRepository dictionaryRepo.DictionaryRepositoryInterface,
	loggerInstance *logger.Logger,
) ISysDictionaryService {
	return &SysDictionaryUseCase{
		sysDictionaryRepository: sysDictionaryRepository,
		Logger:                  loggerInstance,
	}
}

func (s *SysDictionaryUseCase) GetAll() (*[]dictionaryDomain.Dictionary, error) {
	s.Logger.Info("Getting all dictionaries")
	return s.sysDictionaryRepository.GetAll()
}

func (s *SysDictionaryUseCase) GetByID(id int) (*dictionaryDomain.Dictionary, error) {
	s.Logger.Info("Getting dictionary by ID", zap.Int("id", id))
	return s.sysDictionaryRepository.GetByID(id)
}

func (s *SysDictionaryUseCase) Create(newDictionary *dictionaryDomain.Dictionary) (*dictionaryDomain.Dictionary, error) {
	s.Logger.Info("Creating new dictionary", zap.String("Name", newDictionary.Name))
	return s.sysDictionaryRepository.Create(newDictionary)
}

func (s *SysDictionaryUseCase) Delete(id int) error {
	s.Logger.Info("Deleting dictionary", zap.Int("id", id))
	return s.sysDictionaryRepository.Delete(id)
}

func (s *SysDictionaryUseCase) Update(id int, userMap map[string]interface{}) (*dictionaryDomain.Dictionary, error) {
	s.Logger.Info("Updating dictionary", zap.Int("id", id))
	return s.sysDictionaryRepository.Update(id, userMap)
}

func (s *SysDictionaryUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[dictionaryDomain.Dictionary], error) {
	s.Logger.Info("Searching dictionary with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysDictionaryRepository.SearchPaginated(filters)
}

func (s *SysDictionaryUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching dictionary by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysDictionaryRepository.SearchByProperty(property, searchText)
}

func (s *SysDictionaryUseCase) GetOneByMap(userMap map[string]interface{}) (*dictionaryDomain.Dictionary, error) {
	return s.sysDictionaryRepository.GetOneByMap(userMap)
}

func (s *SysDictionaryUseCase) GetByType(typeText string) (*dictionaryDomain.Dictionary, error) {
	return s.sysDictionaryRepository.GetByType(typeText)
}
