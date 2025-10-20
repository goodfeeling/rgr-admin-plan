package role

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainRole "github.com/gbrayhan/microservices-go/src/domain/sys/role"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysRole struct {
	ID            int64          `gorm:"column:id;primary_key;autoIncrement" json:"id,omitempty"`
	CreatedAt     time.Time      `gorm:"column:created_at"`
	UpdatedAt     time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index"`
	Name          string         `gorm:"column:name"`
	ParentID      int64          `gorm:"column:parent_id;type:numeric(20,0)"`
	DefaultRouter string         `gorm:"column:default_router"`
	Status        int16          `gorm:"column:status"`
	Order         int64          `gorm:"column:order;type:numeric(10,0)"`
	Label         string         `gorm:"column:label"`
	Description   string         `gorm:"column:description"`
}

var ColumnsRoleMapping = map[string]string{
	"id":            "id",
	"name":          "name",
	"parentId":      "parent_id",
	"defaultRouter": "default_router",
	"description":   "description",
	"order":         "order",
	"email":         "email",
	"status":        "status",
	"label":         "label",
	"createdAt":     "created_at",
	"updatedAt":     "updated_at",
}

func (SysRole) TableName() string {
	return "sys_roles"
}

type ISysRolesRepository interface {
	GetAll(status int) (*[]domainRole.Role, error)
	Create(roleDomain *domainRole.Role) (*domainRole.Role, error)
	GetByID(id int) (*domainRole.Role, error)
	GetByName(name string) (*domainRole.Role, error)
	Update(id int, roleMap map[string]interface{}) (*domainRole.Role, error)
	Delete(id int) error
	SearchPaginated(filters domain.DataFilters) (*domainRole.SearchResultRole, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(roleMap map[string]interface{}) (*domainRole.Role, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewSysRolesRepository(db *gorm.DB, loggerInstance *logger.Logger) ISysRolesRepository {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll(status int) (*[]domainRole.Role, error) {
	var roles []SysRole
	tx := r.DB
	if status != 0 {
		tx = tx.Where("status = ?", status)
	}
	if err := tx.Find(&roles).Error; err != nil {
		r.Logger.Error("Error getting all roles", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all roles", zap.Int("count", len(roles)))
	return ArrayToDomainMapper(&roles), nil
}

func (r *Repository) Create(roleDomain *domainRole.Role) (*domainRole.Role, error) {
	r.Logger.Info("Creating new role", zap.String("Name", roleDomain.Name))
	roleRepository := fromDomainMapper(roleDomain)
	txDb := r.DB.Create(roleRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating role", zap.Error(err), zap.String("Name", roleDomain.Name))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainRole.Role{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainRole.Role{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created role", zap.String("email", roleDomain.Name), zap.Int("id", int(roleRepository.ID)))
	return roleRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainRole.Role, error) {
	var role SysRole
	err := r.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Role not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting role by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainRole.Role{}, err
	}
	r.Logger.Info("Successfully retrieved role by ID", zap.Int("id", id))
	return role.toDomainMapper(), nil
}

func (r *Repository) GetByName(name string) (*domainRole.Role, error) {
	var role SysRole
	err := r.DB.Where("name = ?", name).First(&role).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Role not found", zap.String("name", name))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting role by name", zap.Error(err), zap.String("name", name))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainRole.Role{}, err
	}
	r.Logger.Info("Successfully retrieved role by name", zap.String("name", name))
	return role.toDomainMapper(), nil
}

func (r *Repository) Update(id int, roleMap map[string]interface{}) (*domainRole.Role, error) {
	var roleObj SysRole
	roleObj.ID = int64(id)
	delete(roleMap, "updated_at")
	err := r.DB.Model(&roleObj).Updates(roleMap).Error
	if err != nil {
		r.Logger.Error("Error updating role", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainRole.Role{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &domainRole.Role{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &domainRole.Role{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&roleObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated role", zap.Error(err), zap.Int("id", id))
		return &domainRole.Role{}, err
	}
	r.Logger.Info("Successfully updated role", zap.Int("id", id))
	return roleObj.toDomainMapper(), nil
}

func (r *Repository) Delete(id int) error {
	tx := r.DB.Delete(&SysRole{}, id)
	if tx.Error != nil {
		r.Logger.Error("Error deleting role", zap.Error(tx.Error), zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Role not found for deletion", zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted role", zap.Int("id", id))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domainRole.SearchResultRole, error) {
	query := r.DB.Model(&SysRole{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsRoleMapping[field]
					if column != "" {
						query = query.Where(column+" ILIKE ?", "%"+value+"%")
					}
				}
			}
		}
	}

	// Apply exact matches
	for field, values := range filters.Matches {
		if len(values) > 0 {
			column := ColumnsRoleMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsRoleMapping[dateFilter.Field]
		if column != "" {
			if dateFilter.Start != nil {
				query = query.Where(column+" >= ?", dateFilter.Start)
			}
			if dateFilter.End != nil {
				query = query.Where(column+" <= ?", dateFilter.End)
			}
		}
	}

	// Apply sorting
	if len(filters.SortBy) > 0 && filters.SortDirection.IsValid() {
		for _, sortField := range filters.SortBy {
			column := ColumnsRoleMapping[sortField]
			if column != "" {
				query = query.Order(column + " " + string(filters.SortDirection))
			}
		}
	}

	// Count total records
	var total int64
	clonedQuery := query
	clonedQuery.Count(&total)

	// Apply pagination
	if filters.Page < 1 {
		filters.Page = 1
	}
	if filters.PageSize < 1 {
		filters.PageSize = 10
	}
	offset := (filters.Page - 1) * filters.PageSize

	var roles []SysRole
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&roles).Error; err != nil {
		r.Logger.Error("Error searching roles", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domainRole.SearchResultRole{
		Data:       ArrayToDomainMapper(&roles),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched roles",
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) SearchByProperty(property string, searchText string) (*[]string, error) {
	column := ColumnsRoleMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysRole{}).
		Distinct(column).
		Where(column+" ILIKE ?", "%"+searchText+"%").
		Limit(20).
		Pluck(column, &coincidences).Error; err != nil {
		r.Logger.Error("Error searching by property", zap.Error(err), zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	r.Logger.Info("Successfully searched by property",
		zap.String("property", property),
		zap.Int("results", len(coincidences)))

	return &coincidences, nil
}

func (u *SysRole) toDomainMapper() *domainRole.Role {
	return &domainRole.Role{
		ID:            u.ID,
		Name:          u.Name,
		ParentID:      u.ParentID,
		Order:         u.Order,
		Label:         u.Label,
		Description:   u.Description,
		Status:        u.Status,
		DefaultRouter: u.DefaultRouter,
		CreatedAt:     domain.CustomTime{Time: u.CreatedAt},
		UpdatedAt:     domain.CustomTime{Time: u.UpdatedAt},
	}
}

func fromDomainMapper(u *domainRole.Role) *SysRole {
	return &SysRole{
		ID:          u.ID,
		Name:        u.Name,
		ParentID:    u.ParentID,
		Order:       u.Order,
		Label:       u.Label,
		Description: u.Description,
		Status:      u.Status,
	}
}

func ArrayToDomainMapper(roles *[]SysRole) *[]domainRole.Role {
	rolesDomain := make([]domainRole.Role, len(*roles))
	for i, role := range *roles {
		rolesDomain[i] = *role.toDomainMapper()
	}
	return &rolesDomain
}

func (r *Repository) GetOneByMap(roleMap map[string]interface{}) (*domainRole.Role, error) {
	var roleRepository SysRole
	tx := r.DB.Limit(1)
	for key, value := range roleMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&roleRepository).Error; err != nil {
		return &domainRole.Role{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return roleRepository.toDomainMapper(), nil
}
