package dictionary_detail

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
)

type DictionaryDetail struct {
	ID              int       `json:"id"`
	Label           string    `json:"label"`
	Value           string    `json:"value"`
	Extend          string    `json:"extend"`
	Status          int16     `json:"status"`
	Sort            int8      `json:"sort"`
	Type            string    `json:"type"`
	SysDictionaryID int64     `json:"sys_dictionary_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type IDictionaryDetailService interface {
	GetAll() (*[]DictionaryDetail, error)
	GetByID(id int) (*DictionaryDetail, error)
	Create(newDictionary *DictionaryDetail) (*DictionaryDetail, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*DictionaryDetail, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[DictionaryDetail], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*DictionaryDetail, error)
}
