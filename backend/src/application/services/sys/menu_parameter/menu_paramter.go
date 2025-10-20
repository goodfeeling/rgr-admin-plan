package menu_parameter

import (
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	menuParameterRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_parameter"

	"github.com/gbrayhan/microservices-go/src/domain"
	menuParameterDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu_parameter"
	"go.uber.org/zap"
)

type IMenuParameterService interface {
	GetAll(menuID int64) (*[]menuParameterDomain.MenuParameter, error)
	GetByID(id int) (*menuParameterDomain.MenuParameter, error)
	Create(newMenuParameter *menuParameterDomain.MenuParameter) (*menuParameterDomain.MenuParameter, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*menuParameterDomain.MenuParameter, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuParameterDomain.MenuParameter], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*menuParameterDomain.MenuParameter, error)
}

type MenuParameterUseCase struct {
	menuParameterRepository menuParameterRepo.MenuParameterRepositoryInterface
	Logger                  *logger.Logger
}

func NewMenuParameterUseCase(menuParameterRepository menuParameterRepo.MenuParameterRepositoryInterface, loggerInstance *logger.Logger) IMenuParameterService {
	return &MenuParameterUseCase{
		menuParameterRepository: menuParameterRepository,
		Logger:                  loggerInstance,
	}
}

func (s *MenuParameterUseCase) GetAll(menuID int64) (*[]menuParameterDomain.MenuParameter, error) {
	s.Logger.Info("Getting all roles")
	return s.menuParameterRepository.GetAll(menuID)
}

func (s *MenuParameterUseCase) GetByID(id int) (*menuParameterDomain.MenuParameter, error) {
	s.Logger.Info("Getting menuParameter by ID", zap.Int("id", id))
	return s.menuParameterRepository.GetByID(id)
}

func (s *MenuParameterUseCase) Create(newMenuParameter *menuParameterDomain.MenuParameter) (*menuParameterDomain.MenuParameter, error) {
	s.Logger.Info("Creating new menuParameter", zap.String("Key", newMenuParameter.Key))
	return s.menuParameterRepository.Create(newMenuParameter)
}

func (s *MenuParameterUseCase) Delete(id int) error {
	s.Logger.Info("Deleting menuParameter", zap.Int("id", id))
	return s.menuParameterRepository.Delete(id)
}

func (s *MenuParameterUseCase) Update(id int, userMap map[string]interface{}) (*menuParameterDomain.MenuParameter, error) {
	s.Logger.Info("Updating menuParameter", zap.Int("id", id))
	return s.menuParameterRepository.Update(id, userMap)
}

func (s *MenuParameterUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuParameterDomain.MenuParameter], error) {
	s.Logger.Info("Searching menuParameter with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.menuParameterRepository.SearchPaginated(filters)
}

func (s *MenuParameterUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching menuParameter by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.menuParameterRepository.SearchByProperty(property, searchText)
}

func (s *MenuParameterUseCase) GetOneByMap(userMap map[string]interface{}) (*menuParameterDomain.MenuParameter, error) {
	return s.menuParameterRepository.GetOneByMap(userMap)
}
