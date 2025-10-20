package casbin_rule

import (
	"fmt"
	"strconv"
	"strings"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// CasbinRule represents the casbin_rule table in the database
type CasbinRule struct {
	ID    int64  `gorm:"primaryKey;column:id;type:numeric(20,0)"` // 主键ID
	PType string `gorm:"column:ptype;type:varchar(100)"`          // 策略类型
	V0    string `gorm:"column:v0;type:varchar(100)"`             // 策略字段 v0
	V1    string `gorm:"column:v1;type:varchar(100)"`             // 策略字段 v1
	V2    string `gorm:"column:v2;type:varchar(100)"`             // 策略字段 v2
	V3    string `gorm:"column:v3;type:varchar(100)"`             // 策略字段 v3
	V4    string `gorm:"column:v4;type:varchar(100)"`             // 策略字段 v4
	V5    string `gorm:"column:v5;type:varchar(100)"`             // 策略字段 v5
}

// TableName returns the name of the database table for this model
func (CasbinRule) TableName() string {
	return "casbin_rule"
}

var ColumnsRoleMapping = map[string]string{
	"id":        "id",
	"createdAt": "created_at",
	"updatedAt": "updated_at",
}

type ICasbinRuleRepository interface {
	Insert(roleId int, UpdateMap map[string]any) error
	GetByRoleId(roleId int) ([]string, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

// GetByRoleId implements ISysRoleMenuRepository.
func (r *Repository) GetByRoleId(roleId int) ([]string, error) {
	var casbinRules []CasbinRule
	err := r.DB.Where(&CasbinRule{V0: strconv.Itoa(roleId)}).Find(&casbinRules).Error
	if err != nil {
		r.Logger.Error("Error getting all casbin_rule", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	var rulePaths []string
	for _, item := range casbinRules {
		rulePaths = append(rulePaths, fmt.Sprintf("%v---%v", item.V1, item.V2))
	}
	return rulePaths, nil
}

// Insert implements ISysRoleMenuRepository.
func (r *Repository) Insert(roleId int, UpdateMap map[string]any) error {
	apiPathsInterface, ok := UpdateMap["apiPaths"]
	if !ok {
		r.Logger.Error("apiPaths not found in update map")
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	apiPathsInterfaceSlice, ok := apiPathsInterface.([]interface{})
	if !ok {
		r.Logger.Error("apiPaths is not an array")
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	casbinMenus := make([]CasbinRule, 0, len(apiPathsInterfaceSlice))
	for _, item := range apiPathsInterfaceSlice {
		apiPathString, ok := item.(string)
		if !ok {
			r.Logger.Error("apiPath is not a string")
			return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		apiPaths := strings.Split(apiPathString, "---")
		roleMenu := CasbinRule{
			PType: "p",
			V0:    strconv.Itoa(roleId),
			V1:    apiPaths[0],
			V2:    apiPaths[1],
		}
		casbinMenus = append(casbinMenus, roleMenu)
	}
	if err := r.DB.Where(&CasbinRule{V0: strconv.Itoa(roleId)}).
		Delete(&CasbinRule{}).Error; err != nil {
		return err
	}
	if err := r.DB.Model(&CasbinRule{}).Create(&casbinMenus).Error; err != nil {
		return err
	}
	return nil
}

func NewCasbinRuleRepository(db *gorm.DB, loggerInstance *logger.Logger) ICasbinRuleRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}
