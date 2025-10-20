package base_menu_parameter

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMenuParameter "github.com/gbrayhan/microservices-go/src/domain/sys/menu_parameter"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysBaseMenuParameter struct {
	CreatedAt     time.Time      `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at" json:"-"`
	SysBaseMenuID int64          `gorm:"column:sys_base_menu_id" json:"menuId"`
	Type          string         `gorm:"column:type" json:"type"`
	Key           string         `gorm:"column:key" json:"key"`
	Value         string         `gorm:"column:value" json:"value"`
	ID            int            `ggorm:"primaryKey;column:id;type:numeric(20,0)"`
}

func (SysBaseMenuParameter) TableName() string {
	return "sys_base_menu_parameters"
}

var ColumnsMenuParameterMapping = map[string]string{
	"id":          "id",
	"path":        "path",
	"menuName":    "menu_name",
	"description": "description",
	"menuGroup":   "menu_group",
	"method":      "method",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

// MenuParameterRepositoryInterface defines the interface for menu repository operations
type MenuParameterRepositoryInterface interface {
	GetAll(menuID int64) (*[]domainMenuParameter.MenuParameter, error)
	Create(menuDomain *domainMenuParameter.MenuParameter) (*domainMenuParameter.MenuParameter, error)
	GetByID(id int) (*domainMenuParameter.MenuParameter, error)
	Update(id int, menuMap map[string]interface{}) (*domainMenuParameter.MenuParameter, error)
	Delete(id int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainMenuParameter.MenuParameter], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(menuMap map[string]interface{}) (*domainMenuParameter.MenuParameter, error)
	GetByIDs(ids []int) (*[]domainMenuParameter.MenuParameter, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewMenuParameterRepository(db *gorm.DB, loggerInstance *logger.Logger) MenuParameterRepositoryInterface {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll(menuID int64) (*[]domainMenuParameter.MenuParameter, error) {
	var menus []SysBaseMenuParameter
	tx := r.DB
	if menuID != 0 {
		tx = tx.Where("sys_base_menu_id = ?", menuID)
	}
	if err := tx.Find(&menus).Error; err != nil {
		r.Logger.Error("Error getting all menus", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all menus", zap.Int("count", len(menus)))
	return ArrayToDomainMapper(&menus), nil
}

func (r *Repository) Create(menuDomain *domainMenuParameter.MenuParameter) (*domainMenuParameter.MenuParameter, error) {
	r.Logger.Info("Creating new menu", zap.String("Key", menuDomain.Key))
	menuRepository := fromDomainMapper(menuDomain)
	txDb := r.DB.Create(menuRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating menu", zap.Error(err), zap.String("Key", menuDomain.Key))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainMenuParameter.MenuParameter{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainMenuParameter.MenuParameter{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created menu", zap.String("Key", menuDomain.Key), zap.Int("id", int(menuRepository.ID)))
	return menuRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainMenuParameter.MenuParameter, error) {
	var menu SysBaseMenuParameter
	err := r.DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("MenuParameter not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting menu by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainMenuParameter.MenuParameter{}, err
	}
	r.Logger.Info("Successfully retrieved menu by ID", zap.Int("id", id))
	return menu.toDomainMapper(), nil
}

func (r *Repository) Update(id int, menuMap map[string]interface{}) (*domainMenuParameter.MenuParameter, error) {
	var menuObj SysBaseMenuParameter
	menuObj.ID = id
	delete(menuMap, "updated_at")
	err := r.DB.Model(&menuObj).
		Select("key", "type", "value").
		Updates(menuMap).Error
	if err != nil {
		r.Logger.Error("Error updating menu", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return nil, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return nil, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&menuObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated menu", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	r.Logger.Info("Successfully updated menu", zap.Int("id", id))
	return menuObj.toDomainMapper(), nil
}

func (r *Repository) Delete(id int) error {
	tx := r.DB.Delete(&SysBaseMenuParameter{}, id)
	if tx.Error != nil {
		r.Logger.Error("Error deleting menu", zap.Error(tx.Error), zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("MenuParameter not found for deletion", zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted menu", zap.Int("id", id))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainMenuParameter.MenuParameter], error) {
	query := r.DB.Model(&SysBaseMenuParameter{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsMenuParameterMapping[field]
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
			column := ColumnsMenuParameterMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsMenuParameterMapping[dateFilter.Field]
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
			column := ColumnsMenuParameterMapping[sortField]
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

	var menus []SysBaseMenuParameter
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&menus).Error; err != nil {
		r.Logger.Error("Error searching menus", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainMenuParameter.MenuParameter]{
		Data:       ArrayToDomainMapper(&menus),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched menus",
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) SearchByProperty(property string, searchText string) (*[]string, error) {
	column := ColumnsMenuParameterMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysBaseMenuParameter{}).
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

func (r *Repository) GetByIDs(ids []int) (*[]domainMenuParameter.MenuParameter, error) {
	var menus []SysBaseMenuParameter
	if err := r.DB.Where("id in (?)", ids).Find(&menus).Error; err != nil {
		r.Logger.Error("Error getting all menus", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all menus", zap.Int("count", len(menus)))
	return ArrayToDomainMapper(&menus), nil
}

func (u *SysBaseMenuParameter) toDomainMapper() *domainMenuParameter.MenuParameter {
	return &domainMenuParameter.MenuParameter{
		ID:            u.ID,
		Key:           u.Key,
		SysBaseMenuID: u.SysBaseMenuID,
		Type:          u.Type,
		Value:         u.Value,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func fromDomainMapper(u *domainMenuParameter.MenuParameter) *SysBaseMenuParameter {
	return &SysBaseMenuParameter{
		ID:            u.ID,
		Key:           u.Key,
		SysBaseMenuID: u.SysBaseMenuID,
		Type:          u.Type,
		Value:         u.Value,
	}
}

func ArrayToDomainMapper(menus *[]SysBaseMenuParameter) *[]domainMenuParameter.MenuParameter {
	menusDomain := make([]domainMenuParameter.MenuParameter, len(*menus))
	for i, menu := range *menus {
		menusDomain[i] = *menu.toDomainMapper()
	}
	return &menusDomain
}

func (r *Repository) GetOneByMap(menuMap map[string]interface{}) (*domainMenuParameter.MenuParameter, error) {
	var menuRepository SysBaseMenuParameter
	tx := r.DB.Limit(1)
	for key, value := range menuMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&menuRepository).Error; err != nil {
		return &domainMenuParameter.MenuParameter{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return menuRepository.toDomainMapper(), nil
}
