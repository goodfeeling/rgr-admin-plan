package files

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type SysFiles struct {
	ID             int64     `json:"id"`
	FileName       string    `json:"file_name"`
	FileMD5        string    `json:"file_md5"`
	FilePath       string    `json:"file_path"`
	FileUrl        string    `json:"file_url"`
	StorageEngine  string    `json:"storage_engine"`
	FileOriginName string    `json:"file_origin_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
type STSTokenCache struct {
	AccessKeyId     string    `json:"access_key_id"`
	AccessKeySecret string    `json:"access_key_secret"`
	SecurityToken   string    `json:"security_token"`
	Expiration      time.Time `json:"expiration"`
	BucketName      string    `json:"bucket_name"`
	Region          string    `json:"region"`
	CreatedAt       time.Time `json:"created_at"`
}

type RefreshToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    string    `json:"user_id"` // 可选，用于关联用户
	CreatedAt time.Time `json:"created_at"`
}
type ISysFilesService interface {
	Create(data *SysFiles) (*SysFiles, error)
	GetAll() (*[]SysFiles, error)
	GetByID(id int) (*SysFiles, error)
	Delete(ids []int64) error
	Update(id int, userMap map[string]interface{}) (*SysFiles, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[SysFiles], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*SysFiles, error)
}
