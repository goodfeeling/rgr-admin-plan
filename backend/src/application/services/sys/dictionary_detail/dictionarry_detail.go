package dictionary_detail

import (
	"fmt"

	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	dictionaryRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary_detail"

	"github.com/gbrayhan/microservices-go/src/domain"
	dictionaryDomain "github.com/gbrayhan/microservices-go/src/domain/sys/dictionary_detail"
	"go.uber.org/zap"
)

type ISysDictionaryService interface {
	GetAll() (*[]dictionaryDomain.DictionaryDetail, error)
	GetByID(id int) (*dictionaryDomain.DictionaryDetail, error)
	Create(newDictionary *dictionaryDomain.DictionaryDetail) (*dictionaryDomain.DictionaryDetail, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*dictionaryDomain.DictionaryDetail, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[dictionaryDomain.DictionaryDetail], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*dictionaryDomain.DictionaryDetail, error)
}

type SysDictionaryUseCase struct {
	sysDictionaryRepository dictionaryRepo.DictionaryRepositoryInterface
	Logger                  *logger.Logger
}

func NewSysDictionaryUseCase(sysDictionaryRepository dictionaryRepo.DictionaryRepositoryInterface, loggerInstance *logger.Logger) ISysDictionaryService {
	return &SysDictionaryUseCase{
		sysDictionaryRepository: sysDictionaryRepository,
		Logger:                  loggerInstance,
	}
}

func (s *SysDictionaryUseCase) GetAll() (*[]dictionaryDomain.DictionaryDetail, error) {
	s.Logger.Info("Getting all dictionary_detail")
	return s.sysDictionaryRepository.GetAll()
}

func (s *SysDictionaryUseCase) GetByID(id int) (*dictionaryDomain.DictionaryDetail, error) {
	s.Logger.Info("Getting dictionary_detailby ID", zap.Int("id", id))
	return s.sysDictionaryRepository.GetByID(id)
}

func (s *SysDictionaryUseCase) Create(newDictionary *dictionaryDomain.DictionaryDetail) (*dictionaryDomain.DictionaryDetail, error) {
	s.Logger.Info("Creating new dictionary_detail", zap.String("Label", newDictionary.Label))
	return s.sysDictionaryRepository.Create(newDictionary)
}

func (s *SysDictionaryUseCase) Delete(ids []int) error {
	s.Logger.Info("Deleting dictionary_detail", zap.String("ids", fmt.Sprintf("%v", ids)))
	return s.sysDictionaryRepository.Delete(ids)
}

func (s *SysDictionaryUseCase) Update(id int, userMap map[string]interface{}) (*dictionaryDomain.DictionaryDetail, error) {
	s.Logger.Info("Updating dictionary", zap.Int("id", id))
	return s.sysDictionaryRepository.Update(id, userMap)
}

func (s *SysDictionaryUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[dictionaryDomain.DictionaryDetail], error) {
	s.Logger.Info("Searching dictionary_detail with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysDictionaryRepository.SearchPaginated(filters)
}

func (s *SysDictionaryUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching dictionary_detail by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysDictionaryRepository.SearchByProperty(property, searchText)
}

func (s *SysDictionaryUseCase) GetOneByMap(userMap map[string]interface{}) (*dictionaryDomain.DictionaryDetail, error) {
	return s.sysDictionaryRepository.GetOneByMap(userMap)
}
