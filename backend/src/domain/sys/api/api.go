package api

import (
	"bytes"
	"mime/multipart"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	"github.com/gin-gonic/gin"
)

type Api struct {
	ID          int       `json:"id"`
	Path        string    `json:"path"`
	ApiGroup    string    `json:"api_group"`
	Method      string    `json:"method"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type IApiService interface {
	GetAll() (*[]Api, error)
	GetByID(id int) (*Api, error)
	Create(newApi *Api) (*Api, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*Api, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[Api], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*Api, error)
	GetApisGroup(path string) (*[]GroupApiItem, error)
	SynchronizeRouterToApi(router gin.RoutesInfo) (*int, error)
	GenerateTemplate() (*bytes.Buffer, error)
	Export() (*bytes.Buffer, error)
	Import(src multipart.File) (*[]Api, *int, *int, error)
}

type GroupApiItem struct {
	GroupName       string          `json:"title"`
	GroupKey        string          `json:"key"`
	DisableCheckbox bool            `json:"disableCheckbox"`
	Children        []*GroupApiItem `json:"children"`
}
