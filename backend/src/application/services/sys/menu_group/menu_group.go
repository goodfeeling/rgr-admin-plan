package menu_group

import (
	"fmt"

	"github.com/gbrayhan/microservices-go/src/domain"
	menuGroupDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu_group"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	menuGroupRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_group"
	"go.uber.org/zap"
)

type ISysMenuGroupService interface {
	GetAll() (*[]menuGroupDomain.MenuGroup, error)
	GetByID(id int) (*menuGroupDomain.MenuGroup, error)
	Create(newMenuGroup *menuGroupDomain.MenuGroup) (*menuGroupDomain.MenuGroup, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*menuGroupDomain.MenuGroup, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuGroupDomain.MenuGroup], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*menuGroupDomain.MenuGroup, error)
}

type SysMenuGroupUseCase struct {
	sysMenuGroupRepository menuGroupRepo.MenuGroupRepositoryInterface
	Logger                 *logger.Logger
}

func NewSysMenuGroupUseCase(sysMenuGroupRepository menuGroupRepo.MenuGroupRepositoryInterface, loggerInstance *logger.Logger) ISysMenuGroupService {
	return &SysMenuGroupUseCase{
		sysMenuGroupRepository: sysMenuGroupRepository,
		Logger:                 loggerInstance,
	}
}

func (s *SysMenuGroupUseCase) GetAll() (*[]menuGroupDomain.MenuGroup, error) {
	s.Logger.Info("Getting all roles")
	return s.sysMenuGroupRepository.GetAll()
}

func (s *SysMenuGroupUseCase) GetByID(id int) (*menuGroupDomain.MenuGroup, error) {
	s.Logger.Info("Getting menuGroup by ID", zap.Int("id", id))
	return s.sysMenuGroupRepository.GetByID(id)
}

func (s *SysMenuGroupUseCase) Create(newMenuGroup *menuGroupDomain.MenuGroup) (*menuGroupDomain.MenuGroup, error) {
	s.Logger.Info("Creating new menuGroup", zap.String("Name", newMenuGroup.Name))
	return s.sysMenuGroupRepository.Create(newMenuGroup)
}

func (s *SysMenuGroupUseCase) Delete(ids []int) error {
	s.Logger.Info("Deleting menuGroup", zap.String("ids", fmt.Sprintf("%v", ids)))
	return s.sysMenuGroupRepository.Delete(ids)
}

func (s *SysMenuGroupUseCase) Update(id int, userMap map[string]interface{}) (*menuGroupDomain.MenuGroup, error) {
	s.Logger.Info("Updating menuGroup", zap.Int("id", id))
	return s.sysMenuGroupRepository.Update(id, userMap)
}

func (s *SysMenuGroupUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuGroupDomain.MenuGroup], error) {
	s.Logger.Info("Searching menuGroups with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysMenuGroupRepository.SearchPaginated(filters)
}

func (s *SysMenuGroupUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching menuGroup by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysMenuGroupRepository.SearchByProperty(property, searchText)
}

// Get one menuGroup by map
func (s *SysMenuGroupUseCase) GetOneByMap(userMap map[string]interface{}) (*menuGroupDomain.MenuGroup, error) {
	return s.sysMenuGroupRepository.GetOneByMap(userMap)
}
