package user

import (
	"errors"
	"os"

	"github.com/gbrayhan/microservices-go/src/application/event/bus"
	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	jwtBlacklistDomain "github.com/gbrayhan/microservices-go/src/domain/jwt_blacklist"
	userDomain "github.com/gbrayhan/microservices-go/src/domain/user"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	userRoleRepo "github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/user_role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	sharedUtil "github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type IUserUseCase interface {
	GetAll() (*[]userDomain.User, error)
	GetByID(id int) (*userDomain.User, error)
	GetByEmail(email string) (*userDomain.User, error)
	Create(newUser *userDomain.User) (*userDomain.User, error)
	Delete(id int) error
	Update(id int64, userMap map[string]interface{}) (*userDomain.User, error)
	SearchPaginated(filters domain.DataFilters) (*userDomain.SearchResultUser, error)
	SearchByProperty(property string, searchText string) (*[]string, error)
	GetOneByMap(userMap map[string]interface{}) (*userDomain.User, error)
	UserBindRoles(userId int64, updateMap map[string]interface{}) error
	ResetPassword(userId int64) (*userDomain.User, error)
	EditPassword(userId int64, data userDomain.PasswordEditRequest) (*userDomain.User, error)
	ChangePasswordById(userId int64, password string, jwtToken string) (*userDomain.User, error)
}

type UserUseCase struct {
	userRepository         user.UserRepositoryInterface
	userRoleRepository     userRoleRepo.ISysUserRoleRepository
	Logger                 *logger.Logger
	eventBus               bus.EventBus
	jwtBlacklistRepository jwtBlacklistDomain.IJwtBlacklistService
}

func NewUserUseCase(
	userRepository user.UserRepositoryInterface,
	userRoleRepository userRoleRepo.ISysUserRoleRepository,
	eventBus bus.EventBus,
	jwtBlacklistRepository jwtBlacklistDomain.IJwtBlacklistService,
	logger *logger.Logger) IUserUseCase {
	return &UserUseCase{
		userRepository:         userRepository,
		userRoleRepository:     userRoleRepository,
		eventBus:               eventBus,
		Logger:                 logger,
		jwtBlacklistRepository: jwtBlacklistRepository,
	}
}

func (s *UserUseCase) GetAll() (*[]userDomain.User, error) {
	s.Logger.Info("Getting all users")
	return s.userRepository.GetAll()
}

func (s *UserUseCase) GetByID(id int) (*userDomain.User, error) {
	s.Logger.Info("Getting user by ID", zap.Int("id", id))
	return s.userRepository.GetByID(id)
}

func (s *UserUseCase) GetByEmail(email string) (*userDomain.User, error) {
	s.Logger.Info("Getting user by email", zap.String("email", email))
	return s.userRepository.GetByEmail(email)
}

func (s *UserUseCase) Create(newUser *userDomain.User) (*userDomain.User, error) {
	s.Logger.Info("Creating new user", zap.String("email", newUser.Email))
	hash, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("Error hashing password", zap.Error(err))
		return &userDomain.User{}, err
	}
	newUser.HashPassword = string(hash)
	newUser.UUID = uuid.New().String()
	newUser.Status = 1
	return s.userRepository.Create(newUser)
}

func (s *UserUseCase) Delete(id int) error {
	s.Logger.Info("Deleting user", zap.Int("id", id))
	return s.userRepository.Delete(id)
}

func (s *UserUseCase) Update(id int64, userMap map[string]interface{}) (*userDomain.User, error) {
	s.Logger.Info("Updating user", zap.Int64("id", id))

	return s.userRepository.Update(id, userMap)
}

func (s *UserUseCase) SearchPaginated(filters domain.DataFilters) (*userDomain.SearchResultUser, error) {
	s.Logger.Info("Searching users with pagination",
		zap.Int("page", filters.Page),
		zap.Int("pageSize", filters.PageSize))
	return s.userRepository.SearchPaginated(filters)
}

func (s *UserUseCase) SearchByProperty(property string, searchText string) (*[]string, error) {
	s.Logger.Info("Searching users by property",
		zap.String("property", property),
		zap.String("searchText", searchText))
	return s.userRepository.SearchByProperty(property, searchText)
}

func (s *UserUseCase) GetOneByMap(userMap map[string]interface{}) (*userDomain.User, error) {
	return s.userRepository.GetOneByMap(userMap)
}
func (s *UserUseCase) UserBindRoles(userId int64, updateMap map[string]interface{}) error {
	return s.userRoleRepository.Insert(userId, updateMap)
}

func (s *UserUseCase) ResetPassword(userId int64) (*userDomain.User, error) {
	updateMap := make(map[string]interface{})
	password := os.Getenv("RESET_USER_PASSWORD")
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("Error hashing password", zap.Error(err))
		return nil, err
	}
	updateMap["hash_password"] = hash
	return s.userRepository.Update(userId, updateMap)
}

func (s *UserUseCase) EditPassword(userId int64, data userDomain.PasswordEditRequest) (*userDomain.User, error) {
	userInfo, err := s.userRepository.GetByID(int(userId))
	if err != nil {
		s.Logger.Error("Error getting user info", zap.Error(err))
		return nil, err
	}
	isAuthenticated := sharedUtil.CheckPasswordHash(data.OldPassword, userInfo.HashPassword)
	if !isAuthenticated {
		s.Logger.Warn("EditPassword failed: invalid password", zap.String("username", userInfo.UserName))
		return nil, domainErrors.NewAppError(errors.New("old password is incorrect"), domainErrors.NotAuthorized)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(data.NewPasswd), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("Error hashing password", zap.Error(err))
		return nil, err
	}
	updateMap := make(map[string]interface{})
	updateMap["hash_password"] = hash
	return s.userRepository.Update(userId, updateMap)
}

// ChangePasswordById implements IUserUseCase.
func (s *UserUseCase) ChangePasswordById(userId int64, password string, jwtToken string) (*userDomain.User, error) {
	s.Logger.Info("Change password by id", zap.Int64("id", userId))
	userInfo, err := s.userRepository.GetByID(int(userId))
	if err != nil {
		s.Logger.Error("Error getting user info", zap.Error(err))
		return nil, err
	}
	updateMap := make(map[string]interface{})
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.Logger.Error("Error hashing password", zap.Error(err))
		return nil, err
	}
	updateMap["hash_password"] = hash
	res, err := s.userRepository.Update(userInfo.ID, updateMap)
	if err != nil {
		s.Logger.Error("Error updating user info", zap.Error(err))
		return nil, err
	}
	// invalidate jwt
	s.jwtBlacklistRepository.AddToBlacklist(jwtToken)
	return res, nil
}
