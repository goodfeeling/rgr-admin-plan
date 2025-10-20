package dictionary

import (
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainDetail "github.com/gbrayhan/microservices-go/src/domain/sys/dictionary_detail"
)

type Dictionary struct {
	ID             int                              `json:"id"`
	Name           string                           `json:"name"`
	Type           string                           `json:"type"`
	Status         int16                            `json:"status"`
	Desc           string                           `json:"desc"`
	IsGenerateFile int16                            `json:"is_generate_file"`
	CreatedAt      time.Time                        `json:"created_at"`
	UpdatedAt      time.Time                        `json:"updated_at"`
	Details        *[]domainDetail.DictionaryDetail `json:"details"`
}

type IDictionaryService interface {
	GetAll() (*[]Dictionary, error)
	GetByID(id int) (*Dictionary, error)
	Create(newDictionary *Dictionary) (*Dictionary, error)
	Delete(id int) error
	Update(id int, userMap map[string]interface{}) (*Dictionary, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[Dictionary], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*Dictionary, error)
	GetByType(typeText string) (*Dictionary, error)
}
