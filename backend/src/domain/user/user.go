package user

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	roleDomain "github.com/gbrayhan/microservices-go/src/domain/sys/role"
)

type User struct {
	ID            int64             `json:"id"`
	UUID          string            `json:"uuid"`
	UserName      string            `json:"user_name"`
	NickName      string            `json:"nick_name"`
	Email         string            `json:"email"`
	Status        int16             `json:"status"`
	HashPassword  string            `json:"hash_password"`
	HeaderImg     string            `json:"header_img"`
	Phone         string            `json:"phone"`
	OriginSetting string            `json:"origin_setting"`
	Password      string            `json:"password"`
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Roles         []roleDomain.Role `json:"roles"`
}
type SearchResultUser struct {
	Data       *[]User `json:"data"`
	Total      int64   `json:"total"`
	Page       int     `json:"page"`
	PageSize   int     `json:"page_size"`
	TotalPages int     `json:"total_page"`
}
type PasswordEditRequest struct {
	ID          int    `json:"id"`
	OldPassword string `json:"oldPassword"`
	NewPasswd   string `json:"newPassword"`
}

type ChangePasswordRequest struct {
	NewPasswd string `json:"new_password"`
}
type IUserService interface {
	GetAll() (*[]User, error)
	GetByID(id int) (*User, error)
	Create(newUser *User) (*User, error)
	Delete(id int) error
	Update(id int64, userMap map[string]interface{}) (*User, error)
	SearchPaginated(filters domain.DataFilters) (*SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*User, error)
	UserBindRoles(userId int64, updateMap map[string]interface{}) error
	ResetPassword(userId int64) (*User, error)
	EditPassword(userId int64, data PasswordEditRequest) (*User, error)
	ChangePasswordById(userId int64, password string, jwtToken string) (*User, error)
}
