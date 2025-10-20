package base_menu_btn

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMenuBtn "github.com/gbrayhan/microservices-go/src/domain/sys/menu_btn"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysBaseMenuBtn struct {
	ID            int            `gorm:"primaryKey;column:id;type:numeric(20,0)"`
	CreatedAt     time.Time      `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt     time.Time      `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;index:idx_sys_apis_deleted_at" json:"deletedAt,omitempty"`
	Name          string         `gorm:"column:name" json:"name,omitempty"`
	Desc          string         `gorm:"column:desc" json:"desc,omitempty"`
	SysBaseMenuID int64          `gorm:"column:sys_base_menu_id" json:"sysBaseMenuBtnId,omitempty"`
}

func (SysBaseMenuBtn) TableName() string {
	return "sys_base_menu_btns"
}

var ColumnsMenuBtnMapping = map[string]string{
	"id":          "id",
	"path":        "path",
	"menuName":    "menu_name",
	"description": "description",
	"menuGroup":   "menu_group",
	"method":      "method",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

// MenuBtnRepositoryInterface defines the interface for menu repository operations
type MenuBtnRepositoryInterface interface {
	GetAll(menuId int64) (*[]domainMenuBtn.MenuBtn, error)
	Create(menuDomain *domainMenuBtn.MenuBtn) (*domainMenuBtn.MenuBtn, error)
	GetByID(id int) (*domainMenuBtn.MenuBtn, error)
	Update(id int, menuMap map[string]interface{}) (*domainMenuBtn.MenuBtn, error)
	Delete(id int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainMenuBtn.MenuBtn], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(menuMap map[string]interface{}) (*domainMenuBtn.MenuBtn, error)
	GetByIDs(ids []int) (*[]domainMenuBtn.MenuBtn, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewMenuBtnRepository(db *gorm.DB, loggerInstance *logger.Logger) MenuBtnRepositoryInterface {
	return &Repository{
		DB:     db,
		Logger: loggerInstance,
	}
}

func (r *Repository) GetAll(menuId int64) (*[]domainMenuBtn.MenuBtn, error) {
	var menus []SysBaseMenuBtn
	tx := r.DB
	if menuId != 0 {
		tx = tx.Where("sys_base_menu_id = ?", menuId)
	}
	if err := tx.Find(&menus).Error; err != nil {
		r.Logger.Error("Error getting all menus", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all menus", zap.Int("count", len(menus)))
	return ArrayToDomainMapper(&menus), nil
}

func (r *Repository) Create(menuDomain *domainMenuBtn.MenuBtn) (*domainMenuBtn.MenuBtn, error) {
	r.Logger.Info("Creating new menu", zap.String("Name", menuDomain.Name))
	menuRepository := fromDomainMapper(menuDomain)
	txDb := r.DB.Create(menuRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating menu", zap.Error(err), zap.String("Name", menuDomain.Name))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainMenuBtn.MenuBtn{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainMenuBtn.MenuBtn{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created menu", zap.String("Name", menuDomain.Name), zap.Int("id", int(menuRepository.ID)))
	return menuRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainMenuBtn.MenuBtn, error) {
	var menu SysBaseMenuBtn
	err := r.DB.Where("id = ?", id).First(&menu).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("MenuBtn not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting menu by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainMenuBtn.MenuBtn{}, err
	}
	r.Logger.Info("Successfully retrieved menu by ID", zap.Int("id", id))
	return menu.toDomainMapper(), nil
}

func (r *Repository) Update(id int, menuMap map[string]interface{}) (*domainMenuBtn.MenuBtn, error) {
	var menuObj SysBaseMenuBtn
	menuObj.ID = id
	delete(menuMap, "updated_at")
	err := r.DB.Model(&menuObj).
		Select("name", "desc").
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
	tx := r.DB.Delete(&SysBaseMenuBtn{}, id)
	if tx.Error != nil {
		r.Logger.Error("Error deleting menu", zap.Error(tx.Error), zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("MenuBtn not found for deletion", zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted menu", zap.Int("id", id))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainMenuBtn.MenuBtn], error) {
	query := r.DB.Model(&SysBaseMenuBtn{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsMenuBtnMapping[field]
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
			column := ColumnsMenuBtnMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsMenuBtnMapping[dateFilter.Field]
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
			column := ColumnsMenuBtnMapping[sortField]
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

	var menus []SysBaseMenuBtn
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&menus).Error; err != nil {
		r.Logger.Error("Error searching menus", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainMenuBtn.MenuBtn]{
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
	column := ColumnsMenuBtnMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysBaseMenuBtn{}).
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

func (r *Repository) GetByIDs(ids []int) (*[]domainMenuBtn.MenuBtn, error) {
	var menus []SysBaseMenuBtn
	if err := r.DB.Where("id in (?)", ids).Find(&menus).Error; err != nil {
		r.Logger.Error("Error getting all menus", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all menus", zap.Int("count", len(menus)))
	return ArrayToDomainMapper(&menus), nil
}

func (u *SysBaseMenuBtn) toDomainMapper() *domainMenuBtn.MenuBtn {
	return &domainMenuBtn.MenuBtn{
		ID:            u.ID,
		Name:          u.Name,
		Desc:          u.Desc,
		SysBaseMenuID: u.SysBaseMenuID,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func fromDomainMapper(u *domainMenuBtn.MenuBtn) *SysBaseMenuBtn {
	return &SysBaseMenuBtn{
		ID:            u.ID,
		Name:          u.Name,
		Desc:          u.Desc,
		SysBaseMenuID: u.SysBaseMenuID,
		CreatedAt:     u.CreatedAt,
		UpdatedAt:     u.UpdatedAt,
	}
}

func ArrayToDomainMapper(menus *[]SysBaseMenuBtn) *[]domainMenuBtn.MenuBtn {
	menusDomain := make([]domainMenuBtn.MenuBtn, len(*menus))
	for i, menu := range *menus {
		menusDomain[i] = *menu.toDomainMapper()
	}
	return &menusDomain
}

func (r *Repository) GetOneByMap(menuMap map[string]interface{}) (*domainMenuBtn.MenuBtn, error) {
	var menuRepository SysBaseMenuBtn
	tx := r.DB.Limit(1)
	for key, value := range menuMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&menuRepository).Error; err != nil {
		return &domainMenuBtn.MenuBtn{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return menuRepository.toDomainMapper(), nil
}
