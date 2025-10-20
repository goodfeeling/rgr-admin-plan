package menu

import (
	"github.com/gbrayhan/microservices-go/src/domain"
	menuBtnDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu_btn"
	menuParameterDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu_parameter"
)

type Menu struct {
	ID             int                                 `json:"id"`
	MenuLevel      int                                 `json:"menu_level"`
	ParentID       int                                 `json:"parent_id"`
	Path           string                              `json:"path"`
	Name           string                              `json:"name"`
	Hidden         bool                                `json:"hidden"`
	Component      string                              `json:"component"`
	Sort           int8                                `json:"sort"`
	KeepAlive      int16                               `json:"keep_alive"`
	Title          string                              `json:"title"`
	Icon           string                              `json:"icon"`
	MenuGroupId    int                                 `json:"menu_group_id"`
	CreatedAt      domain.CustomTime                   `json:"created_at"`
	UpdatedAt      domain.CustomTime                   `json:"updated_at"`
	Level          []int                               `json:"level"`
	Children       []*Menu                             `json:"children"`
	MenuBtns       []menuBtnDomain.MenuBtn             `json:"menu_btns"`
	MenuParameters []menuParameterDomain.MenuParameter `json:"menu_parameters"`
	BtnSlice       []string                            `json:"btn_slice"`
}

type MenuNode struct {
	ID       string      `json:"value"`
	Name     string      `json:"title"`
	Key      string      `json:"key"`
	Path     []int       `json:"path"`
	Children []*MenuNode `json:"children"`
}

type MenuGroup struct {
	Id    int     `json:"id"`
	Name  string  `json:"name"`
	Path  string  `json:"path"`
	Items []*Menu `json:"items"`
}

type IMenuService interface {
	GetAll(groupId int) ([]*Menu, error)
	GetByID(id int) (*Menu, error)
	Create(newMenu *Menu) (*Menu, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*Menu, error)
	GetOneByMap(userMap map[string]interface{}) (*Menu, error)
	GetUserMenus(roleId int64) ([]*MenuGroup, error)
}
