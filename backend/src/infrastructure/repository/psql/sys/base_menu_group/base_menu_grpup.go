package base_menu_group

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"

	"github.com/gbrayhan/microservices-go/src/domain/constants"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainMenuGroup "github.com/gbrayhan/microservices-go/src/domain/sys/menu_group"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	menuRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/base_menu"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SysBaseMenuGroups struct {
	ID        int                    `gorm:"primaryKey;column:id;type:numeric(20,0)"`
	Name      string                 `gorm:"column:name" json:"name"`
	Sort      int8                   `gorm:"column:sort" json:"sort"`
	Path      string                 `gorm:"column:path" json:"path"`
	Status    int16                  `gorm:"column:status" json:"status"`
	CreatedAt time.Time              `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt time.Time              `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time             `gorm:"column:deleted_at" json:"deletedAt"`
	MenuItems []menuRepo.SysBaseMenu `gorm:"foreignKey:MenuGroupId"`
}

func (SysBaseMenuGroups) TableName() string {
	return "sys_base_menu_groups"
}

var ColumnsMenuGroupMapping = map[string]string{
	"id":          "id",
	"path":        "path",
	"apiName":     "api_name",
	"description": "description",
	"apiGroup":    "api_group",
	"method":      "method",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
	"sort":        "sort",
}

// MenuGroupRepositoryInterface defines the interface for api repository operations
type MenuGroupRepositoryInterface interface {
	GetAll() (*[]domainMenuGroup.MenuGroup, error)
	Create(apiDomain *domainMenuGroup.MenuGroup) (*domainMenuGroup.MenuGroup, error)
	GetByID(id int) (*domainMenuGroup.MenuGroup, error)
	Update(id int, apiMap map[string]interface{}) (*domainMenuGroup.MenuGroup, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainMenuGroup.MenuGroup], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(apiMap map[string]interface{}) (*domainMenuGroup.MenuGroup, error)
	GetByRoleId(roleMenuIds []int, roleId int64) (*[]domainMenuGroup.MenuGroup, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewMenuGroupRepository(db *gorm.DB, loggerInstance *logger.Logger) MenuGroupRepositoryInterface {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll() (*[]domainMenuGroup.MenuGroup, error) {
	var apis []SysBaseMenuGroups
	if err := r.DB.
		Preload("MenuItems").
		Preload("MenuItems.MenuBtns").
		Preload("MenuItems.MenuParameters").Find(&apis).Error; err != nil {
		r.Logger.Error("Error getting all apis", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all apis", zap.Int("count", len(apis)))
	return arrayToDomainMapper(&apis), nil
}

func (r *Repository) GetByRoleId(menuIds []int, roleId int64) (*[]domainMenuGroup.MenuGroup, error) {
	var apis []SysBaseMenuGroups

	db := r.DB.Where("status = ?", constants.StatusEnabled)

	if roleId != 0 {
		db = db.Preload("MenuItems", "id in (?)", menuIds)
	} else {
		db = db.Preload("MenuItems")
	}

	if err := db.
		Preload("MenuItems.MenuBtns").
		Preload("MenuItems.MenuParameters").Order("sort asc").Find(&apis).Error; err != nil {
		r.Logger.Error("Error getting all apis", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all apis", zap.Int("count", len(apis)))
	return arrayToDomainMapper(&apis), nil
}

func (r *Repository) Create(apiDomain *domainMenuGroup.MenuGroup) (*domainMenuGroup.MenuGroup, error) {
	r.Logger.Info("Creating new api", zap.String("name", apiDomain.Name))
	apiRepository := fromDomainMapper(apiDomain)
	txDb := r.DB.Create(apiRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating api", zap.Error(err), zap.String("Name", apiDomain.Name))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainMenuGroup.MenuGroup{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainMenuGroup.MenuGroup{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created api", zap.String("Name", apiDomain.Name), zap.Int("id", int(apiRepository.ID)))
	return apiRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainMenuGroup.MenuGroup, error) {
	var api SysBaseMenuGroups
	err := r.DB.Where("id = ?", id).First(&api).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("MenuGroup not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting api by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainMenuGroup.MenuGroup{}, err
	}
	r.Logger.Info("Successfully retrieved api by ID", zap.Int("id", id))
	return api.toDomainMapper(), nil
}

func (r *Repository) Update(id int, apiMap map[string]interface{}) (*domainMenuGroup.MenuGroup, error) {
	var apiObj SysBaseMenuGroups
	apiObj.ID = id
	delete(apiMap, "updated_at")
	err := r.DB.Model(&apiObj).Updates(apiMap).Error
	if err != nil {
		r.Logger.Error("Error updating api", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainMenuGroup.MenuGroup{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &domainMenuGroup.MenuGroup{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &domainMenuGroup.MenuGroup{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&apiObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated api", zap.Error(err), zap.Int("id", id))
		return &domainMenuGroup.MenuGroup{}, err
	}
	r.Logger.Info("Successfully updated api", zap.Int("id", id))
	return apiObj.toDomainMapper(), nil
}

func (r *Repository) Delete(ids []int) error {
	tx := r.DB.Where("id IN ?", ids).Delete(&SysBaseMenuGroups{})

	if tx.Error != nil {
		r.Logger.Error("Error deleting api", zap.Error(tx.Error), zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("MenuGroup not found for deletion", zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted api", zap.String("ids", fmt.Sprintf("%v", ids)))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainMenuGroup.MenuGroup], error) {
	query := r.DB.Model(&SysBaseMenuGroups{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsMenuGroupMapping[field]
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
			column := ColumnsMenuGroupMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsMenuGroupMapping[dateFilter.Field]
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
			column := ColumnsMenuGroupMapping[sortField]
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

	var apis []SysBaseMenuGroups
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&apis).Error; err != nil {
		r.Logger.Error("Error searching apis", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainMenuGroup.MenuGroup]{
		Data:       arrayToDomainMapper(&apis),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched apis",
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) SearchByProperty(property string, searchText string) (*[]string, error) {
	column := ColumnsMenuGroupMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysBaseMenuGroups{}).
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

func (u *SysBaseMenuGroups) toDomainMapper() *domainMenuGroup.MenuGroup {
	return &domainMenuGroup.MenuGroup{
		ID:        u.ID,
		Name:      u.Name,
		Path:      u.Path,
		Status:    u.Status,
		Sort:      u.Sort,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		MenuItems: menuRepo.ArrayToDomainMapper(&u.MenuItems),
	}
}

func fromDomainMapper(u *domainMenuGroup.MenuGroup) *SysBaseMenuGroups {
	return &SysBaseMenuGroups{
		ID:        u.ID,
		Name:      u.Name,
		Path:      u.Path,
		Sort:      u.Sort,
		Status:    u.Status,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func arrayToDomainMapper(apis *[]SysBaseMenuGroups) *[]domainMenuGroup.MenuGroup {
	apisDomain := make([]domainMenuGroup.MenuGroup, len(*apis))
	for i, api := range *apis {
		apisDomain[i] = *api.toDomainMapper()
	}
	return &apisDomain
}

func (r *Repository) GetOneByMap(apiMap map[string]interface{}) (*domainMenuGroup.MenuGroup, error) {
	var apiRepository SysBaseMenuGroups
	tx := r.DB.Limit(1)
	for key, value := range apiMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&apiRepository).Error; err != nil {
		return &domainMenuGroup.MenuGroup{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return apiRepository.toDomainMapper(), nil
}
