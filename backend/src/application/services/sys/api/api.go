package api

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/gbrayhan/microservices-go/src/domain"
	apiDomain "github.com/gbrayhan/microservices-go/src/domain/sys/api"
	"github.com/gbrayhan/microservices-go/src/infrastructure/lib/excel"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	apiRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/api"
	dictionaryRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/dictionary"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type ISysApiService interface {
	GetAll() (*[]apiDomain.Api, error)
	GetByID(id int) (*apiDomain.Api, error)
	Create(newApi *apiDomain.Api) (*apiDomain.Api, error)
	Delete(ids []int) error
	Update(id int, userMap map[string]interface{}) (*apiDomain.Api, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[apiDomain.Api], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*apiDomain.Api, error)
	GetApisGroup(path string) (*[]apiDomain.GroupApiItem, error)
	SynchronizeRouterToApi(router gin.RoutesInfo) (*int, error)
	GenerateTemplate() (*bytes.Buffer, error)
	Export() (*bytes.Buffer, error)
	Import(src multipart.File) (*[]apiDomain.Api, *int, *int, error)
}

type SysApiUseCase struct {
	sysApiRepository     apiRepo.ApiRepositoryInterface
	dictionaryRepository dictionaryRepo.DictionaryRepositoryInterface
	Logger               *logger.Logger
	ExcelHandler         *excel.ExcelHandler
}

func NewSysApiUseCase(
	sysApiRepository apiRepo.ApiRepositoryInterface,
	dictionaryRepository dictionaryRepo.DictionaryRepositoryInterface,
	loggerInstance *logger.Logger) ISysApiService {
	return &SysApiUseCase{
		sysApiRepository:     sysApiRepository,
		dictionaryRepository: dictionaryRepository,
		Logger:               loggerInstance,
		ExcelHandler:         excel.NewExcelHandler(),
	}
}

func (s *SysApiUseCase) GetAll() (*[]apiDomain.Api, error) {
	s.Logger.Info("Getting all roles")
	return s.sysApiRepository.GetAll("")
}

func (s *SysApiUseCase) GetByID(id int) (*apiDomain.Api, error) {
	s.Logger.Info("Getting api by ID", zap.Int("id", id))
	return s.sysApiRepository.GetByID(id)
}

func (s *SysApiUseCase) Create(newApi *apiDomain.Api) (*apiDomain.Api, error) {
	s.Logger.Info("Creating new api", zap.String("path", newApi.Path))
	return s.sysApiRepository.Create(newApi)
}

func (s *SysApiUseCase) Delete(ids []int) error {
	s.Logger.Info("Deleting api", zap.String("ids", fmt.Sprintf("%v", ids)))
	return s.sysApiRepository.Delete(ids)
}

func (s *SysApiUseCase) Update(id int, userMap map[string]interface{}) (*apiDomain.Api, error) {
	s.Logger.Info("Updating api", zap.Int("id", id))
	return s.sysApiRepository.Update(id, userMap)
}

func (s *SysApiUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[apiDomain.Api], error) {
	s.Logger.Info("Searching apis with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysApiRepository.SearchPaginated(filters)
}

func (s *SysApiUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching api by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysApiRepository.SearchByProperty(property, searchText)
}

// Get one api by map
func (s *SysApiUseCase) GetOneByMap(userMap map[string]interface{}) (*apiDomain.Api, error) {
	return s.sysApiRepository.GetOneByMap(userMap)
}

// GetApisGroup
func (s *SysApiUseCase) GetApisGroup(path string) (*[]apiDomain.GroupApiItem, error) {
	apis, err := s.sysApiRepository.GetAll(path)
	if err != nil {
		return nil, err
	}

	dictionary, err := s.dictionaryRepository.GetByType("api_group")
	if err != nil {
		return nil, err
	}

	var groups []apiDomain.GroupApiItem
	for i, item := range *dictionary.Details {
		groupApis := make([]*apiDomain.GroupApiItem, 0)
		for _, api := range *apis {
			if api.ApiGroup == item.Value {
				groupApis = append(groupApis, &apiDomain.GroupApiItem{
					GroupKey:  fmt.Sprintf("%v---%v", api.Path, api.Method),
					GroupName: api.Description,
					Children:  nil,
				})
			}
		}
		if len(groupApis) > 0 {
			group := apiDomain.GroupApiItem{
				GroupName:       item.Label,
				GroupKey:        fmt.Sprintf("0---%v", i), // root node pattern is 0---index
				DisableCheckbox: len(groupApis) == 0,
				Children:        groupApis,
			}
			groups = append(groups, group)
		}

	}
	return &groups, nil
}

func (c *SysApiUseCase) SynchronizeRouterToApi(routes gin.RoutesInfo) (*int, error) {
	count := 0
	for _, route := range routes {
		if c.shouldSyncRoute(route.Path) {
			apiModel := &apiRepo.SysApi{
				Path:        route.Path,
				Method:      route.Method,
				Description: c.generateDescription(route.Path, route.Method),
				ApiGroup:    "other",
			}

			ok, err := c.sysApiRepository.CreateByCondition(apiModel)
			if err != nil {
				c.Logger.Error("Failed to sync route",
					zap.String("path", route.Path),
					zap.String("method", route.Method),
					zap.Error(err))
				continue
			}
			if ok {
				count++
			}
		}
	}
	return &count, nil
}

// GenerateTemplate implements ISysApiService.
func (s *SysApiUseCase) GenerateTemplate() (*bytes.Buffer, error) {
	headers := []string{"Path", "ApiGroup", "Method", "Description"}
	buffer, err := s.ExcelHandler.CreateApiTemplate(headers, "APIs")
	if err != nil {
		s.Logger.Error("Error creating API template", zap.Error(err))

		return nil, err
	}
	return buffer, nil
}

// Export implements ISysApiService.
func (s *SysApiUseCase) Export() (*bytes.Buffer, error) {
	// 获取所有API数据
	apis, err := s.sysApiRepository.GetAll("")
	if err != nil {
		s.Logger.Error("Error getting APIs for export", zap.Error(err))
		return nil, err
	}

	// 转换为Excel数据格式
	headers := []string{"Path", "ApiGroup", "Method", "Description"}
	var rows [][]string

	for _, api := range *apis {
		row := []string{
			fmt.Sprintf("%d", api.ID),
			api.Path,
			api.ApiGroup,
			api.Method,
			api.Description,
		}
		rows = append(rows, row)
	}
	s.Logger.Info("Exporting APIs", zap.Int("count", len(*apis)))

	excelData := &excel.ExcelData{
		Headers: headers,
		Rows:    rows,
	}

	// 创建Excel文件
	buffer, err := s.ExcelHandler.CreateExcel("APIs", excelData)
	if err != nil {
		s.Logger.Error("Error creating Excel file", zap.Error(err))
		return nil, err
	}
	return buffer, nil
}

// Import implements ISysApiService.
func (s *SysApiUseCase) Import(src multipart.File) (*[]apiDomain.Api, *int, *int, error) {
	// 读取Excel数据
	excelData, err := s.ExcelHandler.ReadExcel(src, "APIs")
	if err != nil {
		s.Logger.Error("Error reading Excel file", zap.Error(err))
		return nil, nil, nil, err
	}

	// 验证表头
	expectedHeaders := []string{"Path", "ApiGroup", "Method", "Description"}
	if len(excelData.Headers) != len(expectedHeaders) {
		s.Logger.Error("Invalid Excel format - header count mismatch")
		return nil, nil, nil, err
	}
	// 解析并创建API对象
	var importedApis []apiDomain.Api
	for rowIndex, row := range excelData.Rows {
		if len(row) < 4 {
			s.Logger.Warn("Skipping row due to insufficient columns", zap.Int("row", rowIndex))
			continue
		}

		api := apiDomain.Api{
			Path:        row[0],
			ApiGroup:    row[1],
			Method:      row[2],
			Description: row[3],
		}

		importedApis = append(importedApis, api)
	}

	// 批量创建或更新API
	var createdCount, updatedCount int
	for _, api := range importedApis {
		// 使用upsert操作替换原有的根据ID判断逻辑
		sysApi := &apiRepo.SysApi{
			Path:        api.Path,
			ApiGroup:    api.ApiGroup,
			Method:      api.Method,
			Description: api.Description,
		}
		isCreated, err := s.sysApiRepository.Upsert(sysApi)
		if err != nil {
			s.Logger.Error("Error upserting API during import", zap.Error(err), zap.String("path", api.Path))
			continue
		}

		if isCreated {
			createdCount++
		} else {
			updatedCount++
		}
	}

	return &importedApis, &createdCount, &updatedCount, nil
}
func (a *SysApiUseCase) shouldSyncRoute(path string) bool {
	// 排除一些系统路由
	excludePaths := []string{"/swagger", "/health"}
	for _, exclude := range excludePaths {
		if strings.HasPrefix(path, exclude) {
			return false
		}
	}
	return true
}

func (a *SysApiUseCase) generateDescription(path, method string) string {
	switch method {
	case "GET":
		if strings.Contains(path, "/:id") {
			return "get one resource"
		}
		return "get all resources"
	case "POST":
		return "create resource"
	case "PUT":
		return "update resource"
	case "DELETE":
		return "delete resource"
	default:
		return "api description"
	}
}
