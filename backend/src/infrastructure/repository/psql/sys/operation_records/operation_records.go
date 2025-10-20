package operation_records

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"

	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	domainOperation "github.com/gbrayhan/microservices-go/src/domain/sys/operation_records"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/utils"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// SysOperationRecord represents the sys_operation_records table structure.
type SysOperationRecord struct {
	ID           int            `gorm:"column:id;primary_key;autoIncrement" json:"id,omitempty"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"createdAt,omitempty"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updatedAt,omitempty"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deletedAt,omitempty"`
	IP           string         `gorm:"column:ip" json:"ip,omitempty"`
	Method       string         `gorm:"column:method" json:"method,omitempty"`
	Path         string         `gorm:"column:path" json:"path,omitempty"`
	Status       int64          `gorm:"column:status" json:"status,omitempty"`
	Latency      int64          `gorm:"column:latency" json:"latency,omitempty"`
	Agent        string         `gorm:"column:agent" json:"agent,omitempty"`
	ErrorMessage string         `gorm:"column:error_message" json:"errorMessage,omitempty"`
	Body         string         `gorm:"column:body" json:"body,omitempty"`
	Resp         string         `gorm:"column:resp" json:"resp,omitempty"`
	UserID       int64          `gorm:"column:user_id" json:"userId,omitempty"`
}

func (*SysOperationRecord) TableName() string {
	return "sys_operation_records"
}

var ColumnsOperationMapping = map[string]string{
	"id":     "id",
	"status": "status",
	"path":   "path",
	"method": "method",
}

// OperationRepositoryInterface defines the interface for api repository operations
type OperationRepositoryInterface interface {
	GetAll() (*[]domainOperation.SysOperationRecord, error)
	Create(apiDomain *domainOperation.SysOperationRecord) (*domainOperation.SysOperationRecord, error)
	GetByID(id int) (*domainOperation.SysOperationRecord, error)
	Update(id int, apiMap map[string]interface{}) (*domainOperation.SysOperationRecord, error)
	Delete(ids []int) error
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainOperation.SysOperationRecord], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(apiMap map[string]interface{}) (*domainOperation.SysOperationRecord, error)
}

type Repository struct {
	DB     *gorm.DB
	Logger *logger.Logger
}

func NewOperationRepository(db *gorm.DB, loggerInstance *logger.Logger) OperationRepositoryInterface {
	return &Repository{DB: db, Logger: loggerInstance}
}

func (r *Repository) GetAll() (*[]domainOperation.SysOperationRecord, error) {
	var apis []SysOperationRecord
	if err := r.DB.Find(&apis).Error; err != nil {
		r.Logger.Error("Error getting all apis", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	r.Logger.Info("Successfully retrieved all apis", zap.Int("count", len(apis)))
	return arrayToDomainMapper(&apis), nil
}

func (r *Repository) Create(apiDomain *domainOperation.SysOperationRecord) (*domainOperation.SysOperationRecord, error) {
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
			return &domainOperation.SysOperationRecord{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			err = domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
			return &domainOperation.SysOperationRecord{}, err
		default:
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	r.Logger.Info("Successfully created api", zap.String("Path", apiDomain.Path), zap.Int("id", int(apiRepository.ID)))
	return apiRepository.toDomainMapper(), err
}

func (r *Repository) GetByID(id int) (*domainOperation.SysOperationRecord, error) {
	var api SysOperationRecord
	err := r.DB.Where("id = ?", id).First(&api).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			r.Logger.Warn("Operation not found", zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.NotFound)
		} else {
			r.Logger.Error("Error getting api by ID", zap.Error(err), zap.Int("id", id))
			err = domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
		return &domainOperation.SysOperationRecord{}, err
	}
	r.Logger.Info("Successfully retrieved api by ID", zap.Int("id", id))
	return api.toDomainMapper(), nil
}

func (r *Repository) Update(id int, apiMap map[string]interface{}) (*domainOperation.SysOperationRecord, error) {
	var apiObj SysOperationRecord
	apiObj.ID = id
	delete(apiMap, "updated_at")
	err := r.DB.Model(&apiObj).
		Select("api_name", "email", "nick_name", "status", "phone", "header_img").
		Updates(apiMap).Error
	if err != nil {
		r.Logger.Error("Error updating api", zap.Error(err), zap.Int("id", id))
		byteErr, _ := json.Marshal(err)
		var newError domainErrors.GormErr
		errUnmarshal := json.Unmarshal(byteErr, &newError)
		if errUnmarshal != nil {
			return &domainOperation.SysOperationRecord{}, errUnmarshal
		}
		switch newError.Number {
		case 1062:
			return &domainOperation.SysOperationRecord{}, domainErrors.NewAppErrorWithType(domainErrors.ResourceAlreadyExists)
		default:
			return &domainOperation.SysOperationRecord{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
		}
	}
	if err := r.DB.Where("id = ?", id).First(&apiObj).Error; err != nil {
		r.Logger.Error("Error retrieving updated api", zap.Error(err), zap.Int("id", id))
		return &domainOperation.SysOperationRecord{}, err
	}
	r.Logger.Info("Successfully updated api", zap.Int("id", id))
	return apiObj.toDomainMapper(), nil
}

func (r *Repository) Delete(ids []int) error {
	tx := r.DB.Where("id IN ?", ids).Delete(&SysOperationRecord{})

	if tx.Error != nil {
		r.Logger.Error("Error deleting records", zap.Error(tx.Error), zap.Ints("ids", ids))
		return domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	if tx.RowsAffected == 0 {
		r.Logger.Warn("No records found for deletion", zap.Ints("ids", ids))
		return domainErrors.NewAppErrorWithType(domainErrors.NotFound)
	}

	r.Logger.Info("Successfully deleted records", zap.Ints("ids", ids))
	return nil
}

func (r *Repository) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[domainOperation.SysOperationRecord], error) {
	query := r.DB.Model(&SysOperationRecord{})

	// Apply like filters
	for field, values := range filters.LikeFilters {
		if len(values) > 0 {
			for _, value := range values {
				if value != "" {
					column := ColumnsOperationMapping[field]
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
			column := ColumnsOperationMapping[field]
			if column != "" {
				query = query.Where(column+" IN ?", values)
			}
		}
	}

	// Apply date range filters
	for _, dateFilter := range filters.DateRangeFilters {
		column := ColumnsOperationMapping[dateFilter.Field]
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
			column := ColumnsOperationMapping[sortField]
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

	var apis []SysOperationRecord
	if err := query.Offset(offset).Limit(filters.PageSize).Find(&apis).Error; err != nil {
		r.Logger.Error("Error searching apis", zap.Error(err))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}

	totalPages := int((total + int64(filters.PageSize) - 1) / int64(filters.PageSize))

	result := &domain.PaginatedResult[domainOperation.SysOperationRecord]{
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
	column := ColumnsOperationMapping[property]
	if column == "" {
		r.Logger.Warn("Invalid property for search", zap.String("property", property))
		return nil, domainErrors.NewAppErrorWithType(domainErrors.ValidationError)
	}

	var coincidences []string
	if err := r.DB.Model(&SysOperationRecord{}).
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

func (u *SysOperationRecord) toDomainMapper() *domainOperation.SysOperationRecord {
	return &domainOperation.SysOperationRecord{
		ID:           u.ID,
		IP:           u.IP,
		Path:         u.Path,
		Method:       u.Method,
		Status:       u.Status,
		Latency:      u.Latency,
		Agent:        u.Agent,
		ErrorMessage: u.ErrorMessage,
		Body:         u.Body,

		CreatedAt: domain.CustomTime{Time: u.CreatedAt},
		UpdatedAt: domain.CustomTime{Time: u.UpdatedAt},
	}
}

func arrayToDomainMapper(apis *[]SysOperationRecord) *[]domainOperation.SysOperationRecord {
	apisDomain := make([]domainOperation.SysOperationRecord, len(*apis))
	for i, api := range *apis {
		apisDomain[i] = *api.toDomainMapper()
	}
	return &apisDomain
}

func (r *Repository) GetOneByMap(apiMap map[string]interface{}) (*domainOperation.SysOperationRecord, error) {
	var apiRepository SysOperationRecord
	tx := r.DB.Limit(1)
	for key, value := range apiMap {
		if !utils.IsZeroValue(value) {
			tx = tx.Where(fmt.Sprintf("%s = ?", key), value)
		}
	}
	if err := tx.Find(&apiRepository).Error; err != nil {
		return &domainOperation.SysOperationRecord{}, domainErrors.NewAppErrorWithType(domainErrors.UnknownError)
	}
	return apiRepository.toDomainMapper(), nil
}

func fromDomainMapper(u *domainOperation.SysOperationRecord) *SysOperationRecord {
	return &SysOperationRecord{
		CreatedAt:    u.CreatedAt.Time,
		IP:           u.IP,
		Method:       u.Method,
		Path:         u.Path,
		Status:       u.Status,
		Latency:      u.Latency,
		Agent:        u.Agent,
		ErrorMessage: u.ErrorMessage,
		Body:         u.Body,
		Resp:         u.Resp,
		UserID:       u.UserID,
	}
}
