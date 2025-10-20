package menu_btn

import (
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	menuBtnRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_btn"

	"github.com/gbrayhan/microservices-go/src/domain"
	menuBtnDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu_btn"
	"go.uber.org/zap"
)

type IMenuBtnService interface {
	GetAll(menuId int64) (*[]menuBtnDomain.MenuBtn, error)
	GetByID(id int) (*menuBtnDomain.MenuBtn, error)
	Create(newMenuBtn *menuBtnDomain.MenuBtn) (*menuBtnDomain.MenuBtn, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*menuBtnDomain.MenuBtn, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuBtnDomain.MenuBtn], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*menuBtnDomain.MenuBtn, error)
}

type MenuBtnUseCase struct {
	sysMenuBtnRepository menuBtnRepo.MenuBtnRepositoryInterface
	Logger               *logger.Logger
}

func NewMenuBtnUseCase(sysMenuBtnRepository menuBtnRepo.MenuBtnRepositoryInterface, loggerInstance *logger.Logger) IMenuBtnService {
	return &MenuBtnUseCase{
		sysMenuBtnRepository: sysMenuBtnRepository,
		Logger:               loggerInstance,
	}
}

func (s *MenuBtnUseCase) GetAll(menuID int64) (*[]menuBtnDomain.MenuBtn, error) {
	s.Logger.Info("Getting all roles")
	return s.sysMenuBtnRepository.GetAll(menuID)
}

func (s *MenuBtnUseCase) GetByID(id int) (*menuBtnDomain.MenuBtn, error) {
	s.Logger.Info("Getting menuBtn by ID", zap.Int("id", id))
	return s.sysMenuBtnRepository.GetByID(id)
}

func (s *MenuBtnUseCase) Create(newMenuBtn *menuBtnDomain.MenuBtn) (*menuBtnDomain.MenuBtn, error) {
	s.Logger.Info("Creating new menuBtn", zap.String("Name", newMenuBtn.Name))
	return s.sysMenuBtnRepository.Create(newMenuBtn)
}

func (s *MenuBtnUseCase) Delete(id int) error {
	s.Logger.Info("Deleting menuBtn", zap.Int("id", id))
	return s.sysMenuBtnRepository.Delete(id)
}

func (s *MenuBtnUseCase) Update(id int, userMap map[string]interface{}) (*menuBtnDomain.MenuBtn, error) {
	s.Logger.Info("Updating menuBtn", zap.Int("id", id))
	return s.sysMenuBtnRepository.Update(id, userMap)
}

func (s *MenuBtnUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuBtnDomain.MenuBtn], error) {
	s.Logger.Info("Searching menuBtn with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysMenuBtnRepository.SearchPaginated(filters)
}

func (s *MenuBtnUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching menuBtn by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysMenuBtnRepository.SearchByProperty(property, searchText)
}

func (s *MenuBtnUseCase) GetOneByMap(userMap map[string]interface{}) (*menuBtnDomain.MenuBtn, error) {
	return s.sysMenuBtnRepository.GetOneByMap(userMap)
}
