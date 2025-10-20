package jwt_blacklist

import (
	"time"

	"gorm.io/gorm"
)

type JwtBlacklist struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt *time.Time     `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt *time.Time     `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at" json:"deletedAt,omitempty"`

	Jwt string `gorm:"column:jwt;type:text;uniqueIndex" json:"jwt"`
}

func (*JwtBlacklist) TableName() string {
	return "jwt_blacklists"
}

type JwtBlacklistRepository interface {
	AddToBlacklist(jwtToken string) error
	IsJwtInBlacklist(token string) (bool, error)
}

type Repository struct {
	DB *gorm.DB
}

func NewUJwtBlacklistRepository(db *gorm.DB) JwtBlacklistRepository {
	return &Repository{DB: db}
}

// AddToBlacklist implements JwtBlacklistRepository.
func (r *Repository) AddToBlacklist(jwtToken string) error {
	result := r.DB.Create(&JwtBlacklist{
		Jwt: jwtToken,
	})
	return result.Error
}

// IsJwtInBlacklist implements JwtBlacklistRepository.
func (r *Repository) IsJwtInBlacklist(jwtToken string) (bool, error) {
	var count int64
	r.DB.Model(&JwtBlacklist{}).Where("jwt = ?", jwtToken).Count(&count)
	return count > 0, nil
}
