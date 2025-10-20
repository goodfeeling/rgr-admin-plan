package menu_btn

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type MenuBtn struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Desc          string    `json:"desc"`
	SysBaseMenuID int64     `json:"sys_base_menu_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type IMenuBtnService interface {
	GetAll(menuID int64) (*[]MenuBtn, error)
	GetByID(id int) (*MenuBtn, error)
	Create(newMenu *MenuBtn) (*MenuBtn, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*MenuBtn, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[MenuBtn], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*MenuBtn, error)
}
