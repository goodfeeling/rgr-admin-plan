package files

import (
	"fmt"

	"github.com/gbrayhan/microservices-go/src/domain"
	filesDomain "github.com/gbrayhan/microservices-go/src/domain/sys/files"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/files"
	"go.uber.org/zap"
)

type ISysFilesService interface {
	Create(data *filesDomain.SysFiles) (*filesDomain.SysFiles, error)
	GetAll() (*[]filesDomain.SysFiles, error)
	GetByID(id int) (*filesDomain.SysFiles, error)
	Delete(ids []int64) error
	Update(id int, userMap map[string]interface{}) (*filesDomain.SysFiles, error)
	SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[filesDomain.SysFiles], error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*filesDomain.SysFiles, error)
}

type SysFilesUseCase struct {
	sysFilesRepository files.ISysFilesRepository
	Logger             *logger.Logger
}

// Create implements ISysFilesService.
func (s *SysFilesUseCase) Create(data *filesDomain.SysFiles) (*filesDomain.SysFiles, error) {
	s.Logger.Info("Getting file by filename", zap.String("filename", data.FileName))
	return s.sysFilesRepository.Create(data)
}

func NewSysFilesUseCase(sysFilesRepository files.ISysFilesRepository, loggerInstance *logger.Logger) ISysFilesService {
	return &SysFilesUseCase{
		sysFilesRepository: sysFilesRepository,
		Logger:             loggerInstance,
	}
}

func (s *SysFilesUseCase) GetAll() (*[]filesDomain.SysFiles, error) {
	s.Logger.Info("Getting all files")
	return s.sysFilesRepository.GetAll()
}

func (s *SysFilesUseCase) GetByID(id int) (*filesDomain.SysFiles, error) {
	s.Logger.Info("Getting file by ID", zap.Int("id", id))
	return s.sysFilesRepository.GetByID(id)
}

func (s *SysFilesUseCase) Delete(ids []int64) error {
	s.Logger.Info("Deleting file", zap.String("ids", fmt.Sprintf("%s", ids)))
	return s.sysFilesRepository.Delete(ids)
}

func (s *SysFilesUseCase) Update(id int, userMap map[string]interface{}) (*filesDomain.SysFiles, error) {
	s.Logger.Info("Updating file", zap.Int("id", id))
	return s.sysFilesRepository.Update(id, userMap)
}

func (s *SysFilesUseCase) SearchPaginated(filters domain.DataFilters) (*domain.PaginatedResult[filesDomain.SysFiles], error) {
	s.Logger.Info("Searching file with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.sysFilesRepository.SearchPaginated(filters)
}

func (s *SysFilesUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching file by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.sysFilesRepository.SearchByProperty(property, searchText)
}

func (s *SysFilesUseCase) GetOneByMap(userMap map[string]interface{}) (*filesDomain.SysFiles, error) {
	return s.sysFilesRepository.GetOneByMap(userMap)
}
