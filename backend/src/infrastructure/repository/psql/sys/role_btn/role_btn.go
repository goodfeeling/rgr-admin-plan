package role_btn

import (
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"gorm.io/gorm"
)

// SysRoleBtn 角色按钮关联表
type SysRoleBtn struct {
	RoleID           int64 `gorm:"column:role_id;type:int8;primaryKey"`
	SysMenuID        int64 `gorm:"column:sys_menu_id;type:int8;primaryKey"`
	SysBaseMenuBtnID int64 `gorm:"column:sys_base_menu_btn_id;type:int8;primaryKey"`
}

// TableName 指定表名
func (SysRoleBtn) TableName() string {
	return "sys_role_btns"
}

type ISysRoleBtnRepository interface {
	Insert(roleId int64, UpdateMap map[string]any) error
	GetByRoleId(roleId int64) ([]*SysRoleBtn, error)
}
type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewRoleBtnRepository(db *gorm.DB, loggerInstance *logger.Logger) ISysRoleBtnRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetByRoleId(roleId int64) ([]*SysRoleBtn, error) {
	var roleBtns []*SysRoleBtn
	err := r.DB.Table(SysRoleBtn{}.TableName()).Where("role_id = ?", roleId).Find(&roleBtns).Error
	if err != nil {
		r.Logger.Error(err.Error())
		return nil, err
	}
	return roleBtns, nil
}

// Insert implements ISysRoleMenuRepository.
func (r *Repository) Insert(roleId int64, UpdateMap map[string]any) error {
	btnIdsInterface, ok := UpdateMap["btnIds"]
	if !ok {
		r.Logger.Error("btnIds not found in update map")
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	menuId, ok := UpdateMap["menuId"].(float64)
	if !ok {
		r.Logger.Error("menuId not found in update map")
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	btnIdsInterfaceSlice, ok := btnIdsInterface.([]interface{})
	if !ok {
		r.Logger.Error("btnIds is not an array")
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	roleBtns := make([]SysRoleBtn, 0, len(btnIdsInterfaceSlice))
	for _, item := range btnIdsInterfaceSlice {
		btnIdFloat64, ok := item.(float64)
		if !ok {
			r.Logger.Error("btnId is not a number")
			return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		roleMenu := SysRoleBtn{
			SysBaseMenuBtnID: int64(btnIdFloat64),
			RoleID:           int64(roleId),
			SysMenuID:        int64(menuId),
		}
		roleBtns = append(roleBtns, roleMenu)
	}
	if err := r.DB.Model(&SysRoleBtn{}).
		Where("role_id = ?", roleId).
		Where("sys_menu_id = ?", menuId).
		Delete(&SysRoleBtn{}).Error; err != nil {
		return err
	}
	if len(roleBtns) <= 0 {
		return nil
	}
	if err := r.DB.Model(&SysRoleBtn{}).Create(&roleBtns).Error; err != nil {
		return err
	}
	return nil
}
