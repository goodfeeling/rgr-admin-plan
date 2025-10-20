package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainApi "github.com/gbrayhan/microservices-go/src/domain/sys/api"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// SysApi represents the sys_apis table in the database.
type SysApi struct {
	ID          int            `gorm:"column:id;primary_key" json:"id,string"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index:idx_sys_apis_deleted_at" json:"deletedAt,omitempty"`
	Path        string         `gorm:"column:path" json:"path,omitempty"`               // api路径
	Description string         `gorm:"column:description" json:"description,omitempty"` // api中文描述
	ApiGroup    string         `gorm:"column:api_group" json:"apiGroup,omitempty"`      // api组
	Method      string         `gorm:"column:method" json:"method,omitempty"`           // 方法
}

func (SysApi) TableName() string {
	return "sys_apis"
}

var ColumnsApiMapping = map[string]string{
	"id":          "id",
	"path":        "path",
	"apiName":     "api_name",
	"description": "description",
	"apiGroup":    "api_group",
	"method":      "method",
	"createdAt":   "created_at",
	"updatedAt":   "updated_at",
}

// ApiRepositoryInterface defines the interface for api repository operations
type ApiRepositoryInterface interface {
	GetAll(path string) (*[]domainApi.Api, error)
	Create(apiDomain *domainApi.Api) (*domainApi.Api, error)
	GetByID(id int) (*domainApi.Api, error)
	Update(id int, apiMap map[string]interface{}) (*domainApi.Api, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainApi.Api], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(apiMap map[string]interface{}) (*domainApi.Api, error)
	Upsert(api *SysApi) (bool, error)
	CreateByCondition(api *SysApi) (bool, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewApiRepository(db *gorm.DB, loggerInstance *logger.Logger) ApiRepositoryInterface {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll(path string) (*[]domainApi.Api, error) {
	var apis []SysApi
	handle := r.DB
	if path != "" {
		handle = handle.Where("path LIKE ?", fmt.Sprintf("%%%s%%", path))
	}
	if err := handle.Find(&apis).Error; err != nil {
		r.Logger.Error("Error getting all apis", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all apis", zap.Int("count", len(apis)))
	return arrayToDomainMapper(&apis), nil
}

func (r *Repository) Create(apiDomain *domainApi.Api) (*domainApi.Api, error) {
	r.Logger.Info("Creating new api", zap.String("path", apiDomain.Path))
	apiRepository := fromDomainMapper(apiDomain)
	txDb := r.DB.Create(apiRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating api", zap.Error(err), zap.String("Path", apiDomain.Path))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainApi.Api{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainApi.Api{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created api", zap.String("Path", apiDomain.Path), zap.Int("id", int(apiRepository.ID)))
	return apiRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainApi.Api, error) {
	var api SysApi
	err := r.DB.Where("id = ?", id).First(&api).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Api not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting api by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainApi.Api{}, err
	}
	r.Logger.Info("Successfully retrieved api by ID", zap.Int("id", id))
	return api.toDomainMapper(), nil
}

func (r *Repository) Update(id int, apiMap map[string]interface{}) (*domainApi.Api, error) {
	var apiObj SysApi
	apiObj.ID = id
	delete(apiMap, "updated_at")
	err := r.DB.Model(&apiObj).Updates(apiMap).Error
	if err != nil {
		r.Logger.Error("Error updating api", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainApi.Api{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &domainApi.Api{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &domainApi.Api{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&apiObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated api", zap.Error(err), zap.Int("id", id))
		return &domainApi.Api{}, err
	}
	r.Logger.Info("Successfully updated api", zap.Int("id", id))
	return apiObj.toDomainMapper(), nil
}

func (r *Repository) Delete(ids []int) error {
	tx := r.DB.Where("id IN ?", ids).Delete(&SysApi{})

	if tx.Error != nil {
		r.Logger.Error("Error deleting api", zap.Error(tx.Error), zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Api not found for deletion", zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted api", zap.String("ids", fmt.Sprintf("%v", ids)))
	return nil
}

// update or insert
func (r *Repository) Upsert(api *SysApi) (bool, error) {
	result := r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "path"},
			{Name: "method"},
		},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"api_group":   api.ApiGroup,
			"description": api.Description,
		}),
	}).Create(api)

	if result.Error != nil {
		r.Logger.Error("Failed to upsert api", zap.Error(result.Error))
		return false, result.Error
	}
	isUpdate := result.RowsAffected == 0
	return !isUpdate, nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainApi.Api], error) {
	query := r.DB.Model(&SysApi{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsApiMapping[field]
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
			column := ColumnsApiMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsApiMapping[dateFilter.Field]
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
			column := ColumnsApiMapping[sortField]
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

	var apis []SysApi
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&apis).Error; err != nil {
		r.Logger.Error("Error searching apis", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainApi.Api]{
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
	column := ColumnsApiMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysApi{}).
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

func (r *Repository) CreateByCondition(api *SysApi) (bool, error) {
	var existingApi SysApi

	// 查找是否已存在
	err := r.DB.Where("path = ? AND method = ?", api.Path, api.Method).
		First(&existingApi).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return false, err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 不存在则创建
		result := r.DB.Create(api)
		if result.Error != nil {
			return false, result.Error
		}
	} else {
		return false, nil
	}

	return true, nil
}

func (u *SysApi) toDomainMapper() *domainApi.Api {
	return &domainApi.Api{
		ID:          u.ID,
		Path:        u.Path,
		ApiGroup:    u.ApiGroup,
		Method:      u.Method,
		Description: u.Description,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func fromDomainMapper(u *domainApi.Api) *SysApi {
	return &SysApi{
		ID:          u.ID,
		Path:        u.Path,
		ApiGroup:    u.ApiGroup,
		Method:      u.Method,
		Description: u.Description,
		CreatedAt:   u.CreatedAt,
		UpdatedAt:   u.UpdatedAt,
	}
}

func arrayToDomainMapper(apis *[]SysApi) *[]domainApi.Api {
	apisDomain := make([]domainApi.Api, len(*apis))
	for i, api := range *apis {
		apisDomain[i] = *api.toDomainMapper()
	}
	return &apisDomain
}

func (r *Repository) GetOneByMap(apiMap map[string]interface{}) (*domainApi.Api, error) {
	var apiRepository SysApi
	tx := r.DB.Limit(1)
	for key, value := range apiMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&apiRepository).Error; err != nil {
		return &domainApi.Api{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return apiRepository.toDomainMapper(), nil
}
