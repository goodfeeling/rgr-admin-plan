package user_role

import (
	"strconv"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysUserRole struct {
	SysUserID int64 `gorm:"column:sys_user_id;primaryKey" json:"sysUserId"`
	SysRoleID int64 `gorm:"column:sys_role_id;primaryKey" json:"sysRoleId"`
}

func (SysUserRole) TableName() string {
	return "sys_user_roles"
}

type ISysUserRoleRepository interface {
	Insert(userId int64, UpdateMap map[string]any) error
	GetByUserId(userId int64) ([]int, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

// GetByRoleId implements ISysUserRoleRepository.
func (r *Repository) GetByUserId(userId int64) ([]int, error) {
	var userRoles []SysUserRole
	err := r.DB.Where("sys_user_id = ?", userId).Find(&userRoles).Error
	if err != nil {
		r.Logger.Error("Error getting all roles", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	var roleIds []int
	for _, userRole := range userRoles {
		roleIds = append(roleIds, int(userRole.SysRoleID))
	}
	return roleIds, nil
}

// Insert implements ISysUserRoleRepository.
func (r *Repository) Insert(userId int64, UpdateMap map[string]any) error {

	roleIdsInterface, ok := UpdateMap["roleIds"]

	if !ok {
		r.Logger.Error("roleIds not found in update map")
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	roleIdsInterfaceSlice, ok := roleIdsInterface.([]interface{})
	if !ok {
		r.Logger.Error("roleIds is not an array")
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	userRoles := make([]SysUserRole, 0, len(roleIdsInterfaceSlice))
	for _, item := range roleIdsInterfaceSlice {
		roleIdString, ok := item.(string)
		if !ok {
			r.Logger.Error("roleId is not a number")
			return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		roleId, err := strconv.ParseUint(roleIdString, 10, 64)
		if err != nil {
			r.Logger.Error("roleId is not a number")
			return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		roleMenu := SysUserRole{
			SysUserID: int64(userId),
			SysRoleID: int64(roleId),
		}
		userRoles = append(userRoles, roleMenu)
	}

	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Logger.Error("Failed to begin transaction", zap.Error(tx.Error))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	if err := tx.Model(&SysUserRole{}).Where("sys_user_id = ?", userId).Delete(&SysUserRole{}).Error; err != nil {
		tx.Rollback()
		r.Logger.Error("Failed to delete existing user roles", zap.Error(err))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	if len(userRoles) != 0 {
		if err := tx.Model(&SysUserRole{}).Create(&userRoles).Error; err != nil {
			tx.Rollback()
			r.Logger.Error("Failed to create new user roles", zap.Error(err))
			return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	// 提交事务
	tx.Commit()
	return nil
}

func NewSysUserRoleRepository(db *gorm.DB, loggerInstance *logger.Logger) ISysUserRoleRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}
