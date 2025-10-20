package operation_records

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type SysOperationRecord struct {
	ID           int
	IP           string
	Method       string
	Path         string
	Status       int64
	Latency      int64
	Agent        string
	ErrorMessage string
	Body         string
	Resp         string
	UserID       int64
	CreatedAt    domain.CustomTime
	UpdatedAt    domain.CustomTime
	DeletedAt    time.Time
}
type ISysOperationRecordService interface {
	GetAll() (*[]SysOperationRecord, error)
	GetByID(id int) (*SysOperationRecord, error)
	Create(newSysOperationRecord *SysOperationRecord) (*SysOperationRecord, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*SysOperationRecord, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[SysOperationRecord], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*SysOperationRecord, error)
}
