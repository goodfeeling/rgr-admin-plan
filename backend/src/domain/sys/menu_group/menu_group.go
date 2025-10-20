package menu_group

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	menuDomain "github.com/gbrayhan/microservices-go/src/domain/sys/menu"
)

type MenuGroup struct {
	ID        int                `json:"id"`
	Name      string             `json:"name"`
	Path      string             `json:"path"`
	Sort      int8               `json:"sort"`
	Status    int16              `json:"status"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	MenuItems *[]menuDomain.Menu `json:"menu_items"`
}

type IMenuGroupService interface {
	GetAll() (*[]MenuGroup, error)
	GetByID(id int) (*MenuGroup, error)
	Create(newMenuGroup *MenuGroup) (*MenuGroup, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*MenuGroup, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[MenuGroup], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*MenuGroup, error)
}
