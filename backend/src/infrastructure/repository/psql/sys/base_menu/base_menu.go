package base_menu

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMenu "github.com/gbrayhan/microservices-go/src/domain/sys/menu"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	menuBtnRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_btn"
	menuParamRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu_parameter"

	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysBaseMenu struct {
	ID             int                                  `gorm:"primaryKey;column:id;type:numeric(20,0)"`
	CreatedAt      time.Time                            `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt      time.Time                            `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt      gorm.DeletedAt                       `gorm:"column:deleted_at;index:idx_sys_menus_deleted_at" json:"deletedAt,omitempty"`
	MenuLevel      int                                  `gorm:"column:menu_level;type:numeric(20,0)"`
	ParentID       int                                  `gorm:"column:parent_id;type:numeric(20,0)"`
	Path           string                               `gorm:"column:path"`
	Name           string                               `gorm:"column:name"`
	Hidden         bool                                 `gorm:"column:hidden"`
	Component      string                               `gorm:"column:component"`
	Sort           int8                                 `gorm:"column:sort"`
	KeepAlive      int16                                `gorm:"column:keep_alive"`
	Title          string                               `gorm:"column:title"`
	Icon           string                               `gorm:"column:icon"`
	MenuGroupId    int                                  `gorm:"column:menu_group_id"`
	MenuBtns       []menuBtnRepo.SysBaseMenuBtn         `gorm:"foreignKey:SysBaseMenuID"`
	MenuParameters []menuParamRepo.SysBaseMenuParameter `gorm:"foreignKey:SysBaseMenuID"`
}

func (SysBaseMenu) TableName() string {
	return "sys_base_menus"
}

var ColumnsMenuMapping = map[string]string{
	"id":          "id",
	"path":        "path",
	"menuName":    "menu_name",
	"description": "description",
	"menuGroup":   "menu_group",
	"method":      "method",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

// MenuRepositoryInterface defines the interface for menu repository operations
type MenuRepositoryInterface interface {
	GetAll(groupId int) (*[]domainMenu.Menu, error)
	Create(menuDomain *domainMenu.Menu) (*domainMenu.Menu, error)
	GetByID(id int) (*domainMenu.Menu, error)
	Update(id int, menuMap map[string]interface{}) (*domainMenu.Menu, error)
	Delete(id int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainMenu.Menu], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(menuMap map[string]interface{}) (*domainMenu.Menu, error)
	GetByIDs(ids []int) (*[]domainMenu.Menu, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewMenuRepository(db *gorm.DB, loggerInstance *logger.Logger) MenuRepositoryInterface {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll(groupId int) (*[]domainMenu.Menu, error) {
	var menus []SysBaseMenu
	tx := r.DB
	if groupId != 0 {
		tx = tx.Where("menu_group_id = ?", groupId)
	}
	if err := tx.Find(&menus).Error; err != nil {
		r.Logger.Error("Error getting all menus", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all menus", zap.Int("count", len(menus)))
	return ArrayToDomainMapper(&menus), nil
}

func (r *Repository) Create(menuDomain *domainMenu.Menu) (*domainMenu.Menu, error) {
	r.Logger.Info("Creating new menu", zap.String("path", menuDomain.Path))
	menuRepository := fromDomainMapper(menuDomain)
	txDb := r.DB.Create(menuRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating menu", zap.Error(err), zap.String("Path", menuDomain.Path))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainMenu.Menu{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainMenu.Menu{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created menu", zap.String("Path", menuDomain.Path), zap.Int("id", int(menuRepository.ID)))
	return menuRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainMenu.Menu, error) {
	var menu SysBaseMenu
	err := r.DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Menu not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting menu by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainMenu.Menu{}, err
	}
	r.Logger.Info("Successfully retrieved menu by ID", zap.Int("id", id))
	return menu.toDomainMapper(), nil
}

func (r *Repository) Update(id int, menuMap map[string]interface{}) (*domainMenu.Menu, error) {
	var menuObj SysBaseMenu
	menuObj.ID = id
	delete(menuMap, "updated_at")
	err := r.DB.Model(&menuObj).
		Select("parent_id", "menu_level", "name", "path", "component", "hidden", "sort", "icon", "title", "keep_alive").
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
	tx := r.DB.Delete(&SysBaseMenu{}, id)
	if tx.Error != nil {
		r.Logger.Error("Error deleting menu", zap.Error(tx.Error), zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Menu not found for deletion", zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted menu", zap.Int("id", id))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainMenu.Menu], error) {
	query := r.DB.Model(&SysBaseMenu{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsMenuMapping[field]
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
			column := ColumnsMenuMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsMenuMapping[dateFilter.Field]
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
			column := ColumnsMenuMapping[sortField]
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

	var menus []SysBaseMenu
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&menus).Error; err != nil {
		r.Logger.Error("Error searching menus", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainMenu.Menu]{
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
	column := ColumnsMenuMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysBaseMenu{}).
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

func (r *Repository) GetByIDs(ids []int) (*[]domainMenu.Menu, error) {
	var menus []SysBaseMenu
	if err := r.DB.Where("id in (?)", ids).Find(&menus).Error; err != nil {
		r.Logger.Error("Error getting all menus", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all menus", zap.Int("count", len(menus)))
	return ArrayToDomainMapper(&menus), nil
}

func (u *SysBaseMenu) toDomainMapper() *domainMenu.Menu {
	return &domainMenu.Menu{
		ID:             u.ID,
		Path:           u.Path,
		Name:           u.Name,
		ParentID:       u.ParentID,
		Hidden:         u.Hidden,
		MenuLevel:      u.MenuLevel,
		KeepAlive:      u.KeepAlive,
		Icon:           u.Icon,
		Title:          u.Title,
		Sort:           u.Sort,
		Component:      u.Component,
		MenuGroupId:    u.MenuGroupId,
		MenuBtns:       *menuBtnRepo.ArrayToDomainMapper(&u.MenuBtns),
		MenuParameters: *menuParamRepo.ArrayToDomainMapper(&u.MenuParameters),
		CreatedAt:      domain.CustomTime{Time: u.CreatedAt},
		UpdatedAt:      domain.CustomTime{Time: u.UpdatedAt},
	}
}

func fromDomainMapper(u *domainMenu.Menu) *SysBaseMenu {
	return &SysBaseMenu{
		ID:          u.ID,
		Path:        u.Path,
		Name:        u.Name,
		ParentID:    u.ParentID,
		Hidden:      u.Hidden,
		MenuLevel:   u.MenuLevel,
		KeepAlive:   u.KeepAlive,
		Icon:        u.Icon,
		Title:       u.Title,
		Sort:        u.Sort,
		Component:   u.Component,
		MenuGroupId: u.MenuGroupId,
	}
}

func ArrayToDomainMapper(menus *[]SysBaseMenu) *[]domainMenu.Menu {
	menusDomain := make([]domainMenu.Menu, len(*menus))
	for i, menu := range *menus {
		menusDomain[i] = *menu.toDomainMapper()
	}
	return &menusDomain
}

func (r *Repository) GetOneByMap(menuMap map[string]interface{}) (*domainMenu.Menu, error) {
	var menuRepository SysBaseMenu
	tx := r.DB.Limit(1)
	for key, value := range menuMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&menuRepository).Error; err != nil {
		return &domainMenu.Menu{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return menuRepository.toDomainMapper(), nil
}
