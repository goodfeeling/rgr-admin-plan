package menu_parameter

import (
	"time"
)

type MenuParameter struct {
	ID            int       `json:"id"`
	SysBaseMenuID int64     `json:"menu_id"`
	Type          string    `json:"type"`
	Key           string    `json:"key"`
	Value         string    `json:"value"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type IMenuParameterService interface {
	GetAll(menuID int64) (*[]MenuParameter, error)
	GetByID(id int) (*MenuParameter, error)
	Create(newMenu *MenuParameter) (*MenuParameter, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*MenuParameter, error)
	GetOneByMap(userMap map[string]interface{}) (*MenuParameter, error)
}
