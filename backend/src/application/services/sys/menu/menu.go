package menu

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	menuDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu"
	menuBtnDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu_btn"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"

	menuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"
	menuGroupRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_group"
	roleBtnRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_btn"
	roleMenuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role_menu"
	userRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"go.uber.org/zap"
)

type ISysMenuService interface {
	GetAll(groupId int) ([]*menuDomain.Menu, error)
	GetByID(id int) (*menuDomain.Menu, error)
	Create(newMenu *menuDomain.Menu) (*menuDomain.Menu, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*menuDomain.Menu, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuDomain.Menu], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*menuDomain.Menu, error)
	GetUserMenus(roleId int64) ([]*menuDomain.MenuGroup, error)
}

type SysMenuUseCase struct {
	sysMenuRepository      menuRepo.MenuRepositoryInterface
	userRepository         userRepo.UserRepositoryInterface
	sysRoleMenuRepository  roleMenuRepo.ISysRoleMenuRepository
	sysMenuGroupRepository menuGroupRepo.MenuGroupRepositoryInterface
	sysRoleBtnRepository   roleBtnRepo.ISysRoleBtnRepository
	Logger                 *logger.Logger
}

func NewSysMenuUseCase(
	sysMenuRepository menuRepo.MenuRepositoryInterface,
	sysRoleMenuRepository roleMenuRepo.ISysRoleMenuRepository,
	userRepository userRepo.UserRepositoryInterface,
	sysMenuGroupRepository menuGroupRepo.MenuGroupRepositoryInterface,
	sysRoleBtnRepository roleBtnRepo.ISysRoleBtnRepository,
	loggerInstance *logger.Logger,
) ISysMenuService {
	return &SysMenuUseCase{
		sysMenuRepository:      sysMenuRepository,
		userRepository:         userRepository,
		sysRoleMenuRepository:  sysRoleMenuRepository,
		sysMenuGroupRepository: sysMenuGroupRepository,
		sysRoleBtnRepository:   sysRoleBtnRepository,
		Logger:                 loggerInstance,
	}
}

func (s *SysMenuUseCase) GetAll(groupId int) ([]*menuDomain.Menu, error) {
	s.Logger.Info("Getting all menus")
	menus, err := s.sysMenuRepository.GetAll(groupId)
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, ""), nil
}

func (s *SysMenuUseCase) GetByID(id int) (*menuDomain.Menu, error) {
	s.Logger.Info("Getting menu by ID", zap.Int("id", id))
	return s.sysMenuRepository.GetByID(id)
}

func (s *SysMenuUseCase) Create(newMenu *menuDomain.Menu) (*menuDomain.Menu, error) {
	s.Logger.Info("Creating new menu", zap.String("path", newMenu.Path))
	return s.sysMenuRepository.Create(newMenu)
}

func (s *SysMenuUseCase) Delete(id int) error {
	s.Logger.Info("Deleting menu", zap.Int("id", id))
	return s.sysMenuRepository.Delete(id)
}

func (s *SysMenuUseCase) Update(id int, userMap map[string]interface{}) (*menuDomain.Menu, error) {
	s.Logger.Info("Updating menu", zap.Int("id", id))
	return s.sysMenuRepository.Update(id, userMap)
}

func (s *SysMenuUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[menuDomain.Menu], error) {
	s.Logger.Info("Searching menus with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysMenuRepository.SearchPaginated(filters)
}

func (s *SysMenuUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching menu by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysMenuRepository.SearchByProperty(property, searchText)
}

func (s *SysMenuUseCase) GetOneByMap(userMap map[string]interface{}) (*menuDomain.Menu, error) {
	return s.sysMenuRepository.GetOneByMap(userMap)
}

// GetUserMenus
func (s *SysMenuUseCase) GetUserMenus(roleId int64) ([]*menuDomain.MenuGroup, error) {
	s.Logger.Info("Getting user menus", zap.Int64("roleId", roleId))
	var roleMenuIds []int
	// role bind menu list
	var roleBtns []*roleBtnRepo.SysRoleBtn
	var err error
	if roleId == 0 { // role setting list
		roleMenuIds = []int{}
	} else { // get user menu
		roleMenuIds, err = s.sysRoleMenuRepository.GetByRoleId(roleId)
		if err != nil {
			return nil, err
		}
		roleBtns, err = s.sysRoleBtnRepository.GetByRoleId(roleId)
		if err != nil {
			return nil, err
		}
	}
	s.Logger.Info("getting role btns ", zap.Int("roleBtnsCount", len(roleBtns)))
	roleBtnMap := make(map[int64][]int64)
	for _, roleBtn := range roleBtns {
		if _, exists := roleBtnMap[roleBtn.SysMenuID]; !exists {
			roleBtnMap[roleBtn.SysMenuID] = make([]int64, 0)
		}
		roleBtnMap[roleBtn.SysMenuID] = append(roleBtnMap[roleBtn.SysMenuID], roleBtn.SysBaseMenuBtnID)
	}

	s.Logger.Info("Getting user menus", zap.Int("menusCount", len(roleMenuIds)))
	groups, err := s.sysMenuGroupRepository.GetByRoleId(roleMenuIds, roleId)
	if err != nil {
		return nil, err
	}

	menuGroup := make([]*menuDomain.MenuGroup, 0)
	for _, group := range *groups {
		filteredMenuItems := make([]menuDomain.Menu, 0)
		for _, menuItem := range *group.MenuItems {
			filteredMenuItem := menuItem
			if btnIds, exists := roleBtnMap[int64(menuItem.ID)]; exists && roleId != 0 {
				// 过滤菜单按钮，只保留角色有权访问的按钮
				filteredButtons := make([]menuBtnDomain.MenuBtn, 0)
				btnSlices := make([]string, 0)
				for _, btn := range menuItem.MenuBtns {
					for _, id := range btnIds {
						if int64(btn.ID) == id {
							filteredButtons = append(filteredButtons, btn)
							btnSlices = append(btnSlices, btn.Name)
							break
						}
					}
				}
				filteredMenuItem.MenuBtns = filteredButtons
				filteredMenuItem.BtnSlice = btnSlices
			}
			filteredMenuItems = append(filteredMenuItems, filteredMenuItem)
		}
		treeData := buildMenuTree(&filteredMenuItems, group.Path)
		if treeData == nil {
			treeData = []*menuDomain.Menu{}
		}
		menuGroup = append(menuGroup, &menuDomain.MenuGroup{
			Id:    group.ID,
			Name:  group.Name,
			Path:  group.Path,
			Items: treeData,
		})
	}

	return menuGroup, nil
}

// buildMenuTree
func buildMenuTree(menus *[]menuDomain.Menu, groupPath string) []*menuDomain.Menu {

	menuMap := make(map[int]*menuDomain.Menu)
	var roots []*menuDomain.Menu

	//  traversal: Establish parent-child relationships.
	for _, item := range *menus {
		node := &menuDomain.Menu{
			ID:             item.ID,
			Path:           item.Path,
			Name:           item.Name,
			ParentID:       item.ParentID,
			Hidden:         item.Hidden,
			MenuLevel:      item.MenuLevel,
			KeepAlive:      item.KeepAlive,
			Icon:           item.Icon,
			Title:          item.Title,
			Sort:           item.Sort,
			Component:      item.Component,
			CreatedAt:      item.CreatedAt,
			UpdatedAt:      item.UpdatedAt,
			Level:          []int{item.ID},
			Children:       []*menuDomain.Menu{},
			MenuBtns:       item.MenuBtns,
			MenuParameters: item.MenuParameters,
			BtnSlice:       item.BtnSlice,
		}
		menuMap[item.ID] = node
	}

	// Second traversal: Establish parent-child relationships.
	for _, item := range *menus {
		node := menuMap[item.ID]
		if item.ParentID == 0 {
			roots = append(roots, node)
		} else if parentNode, exists := menuMap[item.ParentID]; exists {
			// path handle
			node.Level = append(node.Level, parentNode.Level...)

			// api get user menu handle
			if groupPath != "" {
				node.Path = parentNode.Path + "/" + node.Path
			}

			parentNode.Children = append(parentNode.Children, node)
		} else {
			// 父节点不存在，作为孤儿节点加入根节点列表
			node.Level = []int{item.ID}
			roots = append(roots, node)
		}
	}

	if groupPath != "" {
		// 只在最终叶子节点添加 groupPath 前缀，并避免重复拼接
		for _, node := range menuMap {
			if len(node.Children) == 0 {
				// 只添加 groupPath 前缀，不重复拼接
				node.Path = "/" + groupPath + "/" + node.Path
			}
		}
	}

	return roots
}
