package dictionary

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	"github.com/gbrayhan/microservices-go/src/domain/constants"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainDictionary "github.com/gbrayhan/microservices-go/src/domain/sys/dictionary"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	dictionaryDetailRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary_detail"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SysDictionary represents the sys_dictionaries table in the database
type SysDictionary struct {
	ID        int            `gorm:"primaryKey;column:id;type:numeric(20,0)"` // 主键ID
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime:milli"`  // 创建时间
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime:milli"`  // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`                 // 软删除标记

	Name           string `gorm:"column:name;type:varchar(191)"` // 字典名（中）
	Type           string `gorm:"column:type;type:varchar(191)"` // 字典名（英）
	Status         int16  `gorm:"column:status;type:smallint"`   // 状态
	Desc           string `gorm:"column:desc;type:varchar(191)"` // 描述
	IsGenerateFile int16  `gorm:"column:is_generate_file;type:smallint"`
	// 关联关系
	Details []dictionaryDetailRepo.SysDictionaryDetail `gorm:"foreignKey:SysDictionaryID"`
}

// TableName returns the name of the database table for this model
func (SysDictionary) TableName() string {
	return "sys_dictionaries"
}

var ColumnsDictionaryMapping = map[string]string{
	"id":             "id",
	"path":           "path",
	"dictionaryName": "dictionary_name",
	"selectedDictId": "dictionary_group",
	"sort":           "sort",
	"createdAt":      "created_at",
	"updatedAt":      "updated_at",
}

// DictionaryRepositoryInterface defines the interface for dictionary repository operations
type DictionaryRepositoryInterface interface {
	GetAll() (*[]domainDictionary.Dictionary, error)
	Create(dictionaryDomain *domainDictionary.Dictionary) (*domainDictionary.Dictionary, error)
	GetByID(id int) (*domainDictionary.Dictionary, error)
	Update(id int, dictionaryMap map[string]interface{}) (*domainDictionary.Dictionary, error)
	Delete(id int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainDictionary.Dictionary], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(dictionaryMap map[string]interface{}) (*domainDictionary.Dictionary, error)

	GetByType(typeText string) (*domainDictionary.Dictionary, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewDictionaryRepository(db *gorm.DB, loggerInstance *logger.Logger) DictionaryRepositoryInterface {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll() (*[]domainDictionary.Dictionary, error) {
	var dictionaries []SysDictionary
	if err := r.DB.Find(&dictionaries).Error; err != nil {
		r.Logger.Error("Error getting all dictionaries", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all dictionaries", zap.Int("count", len(dictionaries)))
	return ArrayToDomainMapper(&dictionaries), nil
}

func (r *Repository) Create(dictionaryDomain *domainDictionary.Dictionary) (*domainDictionary.Dictionary, error) {
	r.Logger.Info("Creating new dictionary", zap.String("Name", dictionaryDomain.Name))
	dictionaryRepository := fromDomainMapper(dictionaryDomain)
	txDb := r.DB.Create(dictionaryRepository)
	err := txDb.Error
	if err != nil {
		r.Logger.Error("Error creating dictionary", zap.Error(err), zap.String("Name", dictionaryDomain.Name))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainDictionary.Dictionary{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainDictionary.Dictionary{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created dictionary", zap.String("Name", dictionaryDomain.Name), zap.Int("id", int(dictionaryRepository.ID)))
	return dictionaryRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainDictionary.Dictionary, error) {
	var dictionary SysDictionary
	err := r.DB.Where("id = ?", id).First(&dictionary).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Dictionary not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting dictionary by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainDictionary.Dictionary{}, err
	}
	r.Logger.Info("Successfully retrieved dictionary by ID", zap.Int("id", id))
	return dictionary.toDomainMapper(), nil
}

func (r *Repository) Update(id int, dictionaryMap map[string]interface{}) (*domainDictionary.Dictionary, error) {
	var dictionaryObj SysDictionary
	dictionaryObj.ID = id
	delete(dictionaryMap, "updated_at")
	err := r.DB.Model(&dictionaryObj).Updates(dictionaryMap).Error
	if err != nil {
		r.Logger.Error("Error updating dictionary", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainDictionary.Dictionary{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &domainDictionary.Dictionary{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &domainDictionary.Dictionary{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&dictionaryObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated dictionary", zap.Error(err), zap.Int("id", id))
		return &domainDictionary.Dictionary{}, err
	}
	r.Logger.Info("Successfully updated dictionary", zap.Int("id", id))
	return dictionaryObj.toDomainMapper(), nil
}

func (r *Repository) Delete(id int) error {
	tx := r.DB.Delete(&SysDictionary{}, id)
	if tx.Error != nil {
		r.Logger.Error("Error deleting dictionary", zap.Error(tx.Error), zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	if tx.RowsAffected == 0 {
		r.Logger.Warn("Dictionary not found for deletion", zap.Int("id", id))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}
	r.Logger.Info("Successfully deleted dictionary", zap.Int("id", id))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainDictionary.Dictionary], error) {
	query := r.DB.Model(&SysDictionary{})

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

	var dictionaries []SysDictionary
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&dictionaries).Error; err != nil {
		r.Logger.Error("Error searching dictionaries", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainDictionary.Dictionary]{
		Data:       ArrayToDomainMapper(&dictionaries),
		Total:      total,
		Page:       filters.Page,
		PageSize:   filters.PageSize,
		TotalPages: totalPages,
	}

	r.Logger.Info("Successfully searched dictionaries",
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
	if err := r.DB.Model(&SysDictionary{}).
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

func (u *SysDictionary) toDomainMapper() *domainDictionary.Dictionary {
	return &domainDictionary.Dictionary{
		ID:             u.ID,
		Name:           u.Name,
		Desc:           u.Desc,
		Type:           u.Type,
		Status:         u.Status,
		IsGenerateFile: u.IsGenerateFile,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
		Details:        dictionaryDetailRepo.ArrayToDomainMapper(&u.Details),
	}
}

func fromDomainMapper(u *domainDictionary.Dictionary) *SysDictionary {
	return &SysDictionary{
		ID:             u.ID,
		Name:           u.Name,
		Desc:           u.Desc,
		Type:           u.Type,
		Status:         u.Status,
		IsGenerateFile: u.IsGenerateFile,
		CreatedAt:      u.CreatedAt,
		UpdatedAt:      u.UpdatedAt,
	}
}

func ArrayToDomainMapper(dictionaries *[]SysDictionary) *[]domainDictionary.Dictionary {
	dictionariesDomain := make([]domainDictionary.Dictionary, len(*dictionaries))
	for i, dictionary := range *dictionaries {
		dictionariesDomain[i] = *dictionary.toDomainMapper()
	}
	return &dictionariesDomain
}

func (r *Repository) GetOneByMap(dictionaryMap map[string]interface{}) (*domainDictionary.Dictionary, error) {
	var dictionaryRepository SysDictionary
	tx := r.DB.Limit(1)
	for key, value := range dictionaryMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&dictionaryRepository).Error; err != nil {
		return &domainDictionary.Dictionary{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return dictionaryRepository.toDomainMapper(), nil
}

func (r *Repository) GetByType(typeText string) (*domainDictionary.Dictionary, error) {
	var dictionaries SysDictionary
	if err := r.DB.
		Preload("Details", func(db *gorm.DB) *gorm.DB {
			return db.Where("status = ?", constants.StatusEnabled).Order("sort desc")
		}).
		Where("type = ? and status = ?", typeText, constants.StatusEnabled).First(&dictionaries).Error; err != nil {
		r.Logger.Error("Error getting all dictionaries", zap.Error(err))
		return nil, nil
	}
	r.Logger.Info("Successfully retrieved all dictionaries", zap.String("typeText", typeText))
	return dictionaries.toDomainMapper(), nil
}
