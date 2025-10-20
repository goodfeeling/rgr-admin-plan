package dictionary_detail

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainDictionary "github.com/gbrayhan/microservices-go/src/domain/sys/dictionary_detail"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SysDictionaryDetailDetail represents the sys_dictionary_details table in the database
type SysDictionaryDetail struct {
	ID        int            `gorm:"primaryKey;column:id;type:numeric(20,0)"` // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`  // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime:milli"`  // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`                 // 软删除标记

	Label           string `gorm:"column:label;type:varchar(191)"`              // 展示值
	Value           string `gorm:"column:value;type:varchar(191)"`              // 字典值
	Extend          string `gorm:"column:extend;type:varchar(191)"`             // 扩展值
	Status          int16  `gorm:"column:status;type:smallint"`                 // 启用状态
	Type            string `gorm:"column:type;type:varchar(30)`                 // 类型
	Sort            int8   `gorm:"column:sort;type:bigint"`                     // 排序标记
	SysDictionaryID int64  `gorm:"column:sys_dictionary_id;type:numeric(20,0)"` // 关联字典ID
}

// TableName returns the name of the database table for this model
func (SysDictionaryDetail) TableName() string {
	return "sys_dictionary_details"
}

var ColumnsDictionaryMapping = map[string]string{
	"id":             "id",
	"selectedDictId": "sys_dictionary_id",
	"sort":           "sort",
	"updatedAt":      "updated_at",
}

// DictionaryRepositoryInterface defines the interface for dictionary repository operations
type DictionaryRepositoryInterface interface {
	GetAll() (*[]domainDictionary.DictionaryDetail, error)
	Create(dictionaryDomain *domainDictionary.DictionaryDetail) (*domainDictionary.DictionaryDetail, error)
	GetByID(id int) (*domainDictionary.DictionaryDetail, error)
	Update(id int, dictionaryMap map[string]interface{}) (*domainDictionary.DictionaryDetail, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainDictionary.DictionaryDetail], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(dictionaryMap map[string]interface{}) (*domainDictionary.DictionaryDetail, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewDictionaryRepository(db *gorm.DB, loggerInstance *logger.Logger) DictionaryRepositoryInterface {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll() (*[]domainDictionary.DictionaryDetail, error) {
	var dictionarys []SysDictionaryDetail
	if err := r.DB.Find(&dictionarys).Error; err != nil {
		r.Logger.Error("Error getting all dictionarys", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all dictionarys", zap.Int("count", len(dictionarys)))
	return ArrayToDomainMapper(&dictionarys), nil
}

func (r *Repository) Create(dictionaryDomain *domainDictionary.DictionaryDetail) (*domainDictionary.DictionaryDetail, error) {
	r.Logger.Info("Creating new dictionary", zap.String("label", dictionaryDomain.Label))
	dictionaryRepository := fromDomainMapper(dictionaryDomain)
	txDb := r.DB.Create(dictionaryRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating dictionary", zap.Error(err), zap.String("Label", dictionaryDomain.Label))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return nil, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return nil, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created dictionary", zap.String("Label", dictionaryDomain.Label), zap.Int("id", int(dictionaryRepository.ID)))
	return dictionaryRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainDictionary.DictionaryDetail, error) {
	var dictionary SysDictionaryDetail
	err := r.DB.Where("id = ?", id).First(&dictionary).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Dictionary not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting dictionary by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return nil, err
	}
	r.Logger.Info("Successfully retrieved dictionary by ID", zap.Int("id", id))
	return dictionary.toDomainMapper(), nil
}

func (r *Repository) Update(id int, dictionaryMap map[string]interface{}) (*domainDictionary.DictionaryDetail, error) {
	var dictionaryObj SysDictionaryDetail
	dictionaryObj.ID = id
	delete(dictionaryMap, "updated_at")
	err := r.DB.Model(&dictionaryObj).Updates(dictionaryMap).Error
	if err != nil {
		r.Logger.Error("Error updating dictionary", zap.Error(err), zap.Int("id", id))
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
	if err := r.DB.Where("id = ?", id).First(&dictionaryObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated dictionary", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	r.Logger.Info("Successfully updated dictionary", zap.Int("id", id))
	return dictionaryObj.toDomainMapper(), nil
}

func (r *Repository) Delete(ids []int) error {
	tx := r.DB.Where("id IN ?", ids).Delete(&SysDictionaryDetail{})

	if tx.Error != nil {
		r.Logger.Error("Error deleting dictionary details", zap.Error(tx.Error), zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Api not found for deletion", zap.String("ids", fmt.Sprintf("%v", ids)))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted dictionary details", zap.String("ids", fmt.Sprintf("%v", ids)))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainDictionary.DictionaryDetail], error) {
	query := r.DB.Model(&SysDictionaryDetail{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsDictionaryMapping[field]
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
			column := ColumnsDictionaryMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsDictionaryMapping[dateFilter.Field]
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
			column := ColumnsDictionaryMapping[sortField]
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

	var dictionarys []SysDictionaryDetail
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&dictionarys).Error; err != nil {
		r.Logger.Error("Error searching dictionarys", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainDictionary.DictionaryDetail]{
		Data:       ArrayToDomainMapper(&dictionarys),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched dictionarys",
		zap.Int64("total", total),
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))

	return result, nil
}

func (r *Repository) SearchByProperty(property string, searchText string) (*[]string, error) {
	column := ColumnsDictionaryMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysDictionaryDetail{}).
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

func (u *SysDictionaryDetail) toDomainMapper() *domainDictionary.DictionaryDetail {
	return &domainDictionary.DictionaryDetail{
		ID:              u.ID,
		Label:           u.Label,
		Value:           u.Value,
		Extend:          u.Extend,
		Status:          u.Status,
		Sort:            u.Sort,
		Type:            u.Type,
		SysDictionaryID: u.SysDictionaryID,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func fromDomainMapper(u *domainDictionary.DictionaryDetail) *SysDictionaryDetail {
	return &SysDictionaryDetail{
		ID:              u.ID,
		Label:           u.Label,
		Value:           u.Value,
		Extend:          u.Extend,
		Status:          u.Status,
		Sort:            u.Sort,
		Type:            u.Type,
		SysDictionaryID: u.SysDictionaryID,
		CreatedAt:       u.CreatedAt,
		UpdatedAt:       u.UpdatedAt,
	}
}

func ArrayToDomainMapper(dictionarys *[]SysDictionaryDetail) *[]domainDictionary.DictionaryDetail {
	dictionarysDomain := make([]domainDictionary.DictionaryDetail, len(*dictionarys))
	for i, dictionary := range *dictionarys {
		dictionarysDomain[i] = *dictionary.toDomainMapper()
	}
	return &dictionarysDomain
}

func (r *Repository) GetOneByMap(dictionaryMap map[string]interface{}) (*domainDictionary.DictionaryDetail, error) {
	var dictionaryRepository SysDictionaryDetail
	tx := r.DB.Limit(1)
	for key, value := range dictionaryMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&dictionaryRepository).Error; err != nil {
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return dictionaryRepository.toDomainMapper(), nil
}
