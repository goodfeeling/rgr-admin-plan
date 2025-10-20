package role

import (
	"strconv"

	menuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"

	"github.com/gbrayhan/microservices-go/src/domain"
	roleDomain "github.com/gbrayhan/microservices-go/src/domain/sys/role"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	casbinRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/casbin_rule"
	roleRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	roleBtnRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_btn"
	roleMenuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_menu"
	"go.uber.org/zap"
)

type ISysRoleService interface {
	GetAll(status int) ([]*roleDomain.RoleTree, error)
	GetByID(id int) (*roleDomain.Role, error)
	GetByName(name string) (*roleDomain.Role, error)
	Create(newRole *roleDomain.Role) (*roleDomain.Role, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*roleDomain.Role, error)
	SearchPaginated(filters domain.DataFilters) (*roleDomain.SearchResultRole, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*roleDomain.Role, error)
	GetTreeRoles(status int) (*roleDomain.RoleNode, error)

	GetRoleMenuIds(id int64) (map[int][]int, map[int64][]int64, error)
	UpdateRoleMenuIds(id int, updateMap map[string]any) error

	GetApiRuleList(roleId int) ([]string, error)
	BindApiRule(roleId int, updateMap map[string]interface{}) error
	BindRoleMenuBtns(roleId int64, updateMap map[string]interface{}) error
}

type SysRoleUseCase struct {
	sysRoleRepository     roleRepo.ISysRolesRepository
	sysRoleMenuRepository roleMenuRepo.ISysRoleMenuRepository
	sysMenuRepository     menuRepo.MenuRepositoryInterface
	casbinRuleRepo        casbinRepo.ICasbinRuleRepository
	sysRoleBtnRepo        roleBtnRepo.ISysRoleBtnRepository

	Logger *logger.Logger
}

func NewSysRoleUseCase(
	sysRoleRepository roleRepo.ISysRolesRepository,
	sysRoleMenuRepository roleMenuRepo.ISysRoleMenuRepository,
	casbinRuleRepo casbinRepo.ICasbinRuleRepository,
	sysMenuRepository menuRepo.MenuRepositoryInterface,
	sysRoleBtnRepo roleBtnRepo.ISysRoleBtnRepository,
	loggerInstance *logger.Logger) ISysRoleService {
	return &SysRoleUseCase{
		sysRoleRepository:     sysRoleRepository,
		sysRoleMenuRepository: sysRoleMenuRepository,
		sysMenuRepository:     sysMenuRepository,
		sysRoleBtnRepo:        sysRoleBtnRepo,
		casbinRuleRepo:        casbinRuleRepo,
		Logger:                loggerInstance,
	}
}

func (s *SysRoleUseCase) GetAll(status int) ([]*roleDomain.RoleTree, error) {

	s.Logger.Info("Getting all roles")
	roles, err := s.sysRoleRepository.GetAll(0)
	if err != nil {
		return nil, err
	}
	return BuildRoleTree(roles), nil
}

func BuildRoleTree(roles *[]roleDomain.Role) []*roleDomain.RoleTree {
	roleMap := make(map[int64]*roleDomain.RoleTree)
	var roots []*roleDomain.RoleTree

	// First traversal: Create all nodes and put them into the map.
	for _, role := range *roles {

		node := &roleDomain.RoleTree{
			ID:            role.ID,
			Name:          role.Name,
			ParentID:      role.ParentID,
			DefaultRouter: role.DefaultRouter,
			Status:        role.Status,
			Order:         role.Order,
			Label:         role.Label,
			Description:   role.Description,
			CreatedAt:     role.CreatedAt,
			UpdatedAt:     role.UpdatedAt,
			Path:          []int64{role.ID},
			Children:      []*roleDomain.RoleTree{},
		}
		roleMap[role.ID] = node
	}

	// Second traversal: Establish parent-child relationships.
	for _, role := range *roles {
		node := roleMap[role.ID]
		if role.ParentID == 0 {
			roots = append(roots, node)
		} else if parentNode, exists := roleMap[role.ParentID]; exists {
			// path handle
			node.Path = append(node.Path, parentNode.Path...)
			parentNode.Children = append(parentNode.Children, node)
		} else {
			// 父节点不存在，作为孤儿节点加入根节点列表
			roots = append(roots, node)
		}
	}

	return roots
}

func (s *SysRoleUseCase) GetByID(id int) (*roleDomain.Role, error) {
	s.Logger.Info("Getting role by ID", zap.Int("id", id))
	return s.sysRoleRepository.GetByID(id)
}

func (s *SysRoleUseCase) GetByName(name string) (*roleDomain.Role, error) {
	s.Logger.Info("Getting role by name", zap.String("name", name))
	return s.sysRoleRepository.GetByName(name)
}

func (s *SysRoleUseCase) Create(newRole *roleDomain.Role) (*roleDomain.Role, error) {
	s.Logger.Info("Creating new role", zap.String("name", newRole.Name))
	return s.sysRoleRepository.Create(newRole)
}

func (s *SysRoleUseCase) Delete(id int) error {
	s.Logger.Info("Deleting role", zap.Int("id", id))
	return s.sysRoleRepository.Delete(id)
}

func (s *SysRoleUseCase) Update(id int, userMap map[string]interface{}) (*roleDomain.Role, error) {
	s.Logger.Info("Updating role", zap.Int("id", id))
	return s.sysRoleRepository.Update(id, userMap)
}

func (s *SysRoleUseCase) SearchPaginated(filters domain.DataFilters) (*roleDomain.SearchResultRole, error) {
	s.Logger.Info("Searching roles with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysRoleRepository.SearchPaginated(filters)
}

func (s *SysRoleUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching role by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysRoleRepository.SearchByProperty(property, searchText)
}

func (s *SysRoleUseCase) GetOneByMap(userMap map[string]interface{}) (*roleDomain.Role, error) {
	return s.sysRoleRepository.GetOneByMap(userMap)
}

// GetTreeRoles implements ISysRoleService.
func (s *SysRoleUseCase) GetTreeRoles(status int) (*roleDomain.RoleNode, error) {
	roles, err := s.sysRoleRepository.GetAll(status)
	if err != nil {
		return nil, err
	}
	roleMap := make(map[int64]*roleDomain.RoleNode)
	var roots []*roleDomain.RoleNode

	// First traversal: Create all nodes and put them into the map.
	for _, role := range *roles {
		id := strconv.Itoa(int(role.ID))
		node := &roleDomain.RoleNode{
			ID:       id,
			Name:     role.Name,
			Key:      id,
			Path:     []int64{role.ID},
			Children: []*roleDomain.RoleNode{},
		}
		roleMap[role.ID] = node
	}

	// Second traversal: Establish parent-child relationships.
	for _, role := range *roles {
		node := roleMap[role.ID]
		if role.ParentID == 0 {
			roots = append(roots, node)
		} else if parentNode, exists := roleMap[role.ParentID]; exists {
			// path handle
			node.Path = append(node.Path, parentNode.Path...)
			parentNode.Children = append(parentNode.Children, node)
		} else {
			roots = append(roots, node)
		}
	}
	return &roleDomain.RoleNode{
		ID:       "0",
		Name:     "根节点",
		Key:      "0",
		Children: roots,
	}, nil
}

func (s *SysRoleUseCase) GetRoleMenuIds(id int64) (map[int][]int, map[int64][]int64, error) {
	menuIds, err := s.sysRoleMenuRepository.GetByRoleId(id)
	if err != nil {
		return nil, nil, err
	}
	menus, err := s.sysMenuRepository.GetByIDs(menuIds)
	menuGroups := make(map[int][]int, 0)
	for _, v := range *menus {
		menuGroups[v.MenuGroupId] = append(menuGroups[v.MenuGroupId], v.ID)
	}
	menuBtns, err := s.sysRoleBtnRepo.GetByRoleId(id)
	roleBtns := make(map[int64][]int64, 0)
	for _, v := range menuBtns {
		roleBtns[v.SysMenuID] = append(roleBtns[v.SysMenuID], v.SysBaseMenuBtnID)
	}
	return menuGroups, roleBtns, nil
}

func (s *SysRoleUseCase) UpdateRoleMenuIds(id int, updateMap map[string]any) error {
	return s.sysRoleMenuRepository.Insert(id, updateMap)
}
func (s *SysRoleUseCase) GetApiRuleList(roleId int) ([]string, error) {
	return s.casbinRuleRepo.GetByRoleId(roleId)

}
func (s *SysRoleUseCase) BindApiRule(roleId int, updateMap map[string]interface{}) error {
	return s.casbinRuleRepo.Insert(roleId, updateMap)
}

func (s *SysRoleUseCase) BindRoleMenuBtns(roleId int64, updateMap map[string]interface{}) error {
	return s.sysRoleBtnRepo.Insert(roleId, updateMap)
}
