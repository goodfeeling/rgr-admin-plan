package auth

import (
	"context"
	"errors"
	"os"
	"time"

	"github.com/gbrayhan/microservices-go/src/domain"
	domainErrors "github.com/gbrayhan/microservices-go/src/domain/errors"
	jwtBlacklistDomain "github.com/gbrayhan/microservices-go/src/domain/jwt_blacklist"
	domainRole "github.com/gbrayhan/microservices-go/src/domain/sys/role"
	domainUser "github.com/gbrayhan/microservices-go/src/domain/user"
	captchaLib "github.com/gbrayhan/microservices-go/src/infrastructure/lib/captcha"
	logger "github.com/gbrayhan/microservices-go/src/infrastructure/lib/logger"
	ws "github.com/gbrayhan/microservices-go/src/infrastructure/lib/websocket"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/sys/role"
	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/user"
	"github.com/gbrayhan/microservices-go/src/infrastructure/security"
	sharedUtil "github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type IAuthUseCase interface {
	Login(username, password, captchaId, CaptchaAnswer string, ctx *gin.Context) (*domainUser.User, *AuthTokens, *domainRole.Role, error)
	Logout(jwtToken string) (*domain.CommonResponse[string], error)
	Register(user RegisterUser) (*domain.CommonResponse[SecurityRegisterUser], error)
	AccessTokenByRefreshToken(refreshToken string) (*domainUser.User, *AuthTokens, error)
	SwitchRole(userId int, roleId int64) (*domainUser.User, *AuthTokens, *domainRole.Role, error)
}

type AuthUseCase struct {
	UserRepository         user.UserRepositoryInterface
	RoleRepository         role.ISysRolesRepository
	JWTService             security.IJWTService
	Logger                 *logger.Logger
	jwtBlacklistRepository jwtBlacklistDomain.IJwtBlacklistService
	RedisClient            *redis.Client
	sessionManager         *ws.SessionManager
	captchaHandler         *captchaLib.Captcha
}

func NewAuthUseCase(
	userRepository user.UserRepositoryInterface,
	RoleRepository role.ISysRolesRepository,
	jwtService security.IJWTService,
	loggerInstance *logger.Logger,
	jwtBlacklistRepository jwtBlacklistDomain.IJwtBlacklistService,
	RedisClient *redis.Client,
	sessionManager *ws.SessionManager,
	captchaHandler *captchaLib.Captcha,
) IAuthUseCase {
	return &AuthUseCase{
		UserRepository:         userRepository,
		RoleRepository:         RoleRepository,
		JWTService:             jwtService,
		Logger:                 loggerInstance,
		jwtBlacklistRepository: jwtBlacklistRepository,
		RedisClient:            RedisClient,
		sessionManager:         sessionManager,
		captchaHandler:         captchaHandler,
	}
}

type AuthTokens struct {
	AccessToken               string
	RefreshToken              string
	ExpirationAccessDateTime  time.Time
	ExpirationRefreshDateTime time.Time
}

func (s *AuthUseCase) SwitchRole(userId int, roleId int64) (*domainUser.User, *AuthTokens, *domainRole.Role, error) {
	s.Logger.Info("User switch attempt", zap.Int("userId", userId))
	user, err := s.UserRepository.GetByID(int(userId))
	if err != nil {
		s.Logger.Error("Error getting user for switch", zap.Error(err), zap.Int("userId", userId))
		return nil, nil, nil, err
	}
	if user.ID == 0 {
		s.Logger.Warn("Login failed: user not found", zap.Int("userId", userId))
		return nil, nil, nil, domainErrors.NewAppError(errors.New("user don't no found"), domainErrors.NotAuthorized)
	}
	role, err := s.RoleRepository.GetByID(int(roleId))
	if err != nil {
		s.Logger.Error("Error getting role for switch", zap.Error(err), zap.Int("roleId", int(roleId)))
		return nil, nil, nil, err
	}
	accessTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "access")
	if err != nil {
		s.Logger.Error("Error generating access token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, nil, err
	}
	refreshTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "refresh")
	if err != nil {
		s.Logger.Error("Error generating refresh token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, nil, err
	}

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		RefreshToken:              refreshTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		ExpirationRefreshDateTime: refreshTokenClaims.ExpirationTime,
	}

	s.Logger.Info("User login successful", zap.Int("userId", userId))
	return user, authTokens, role, nil
}

func (s *AuthUseCase) Login(username, password, captchaId, CaptchaAnswer string, ginCtx *gin.Context) (*domainUser.User, *AuthTokens, *domainRole.Role, error) {
	isValid := s.captchaHandler.Verify(captchaId, CaptchaAnswer)
	if !isValid {
		return nil, nil, nil, domainErrors.NewAppError(errors.New("invalid captcha"), domainErrors.CaptchaError)
	}
	user, err := s.UserRepository.GetByUsername(username)
	if err != nil {
		s.Logger.Error("Error getting user for login", zap.Error(err), zap.String("username", username))
		return nil, nil, nil, err
	}
	if user.ID == 0 {
		s.Logger.Warn("Login failed: user not found", zap.String("username", username))
		return nil, nil, nil, domainErrors.NewAppError(errors.New("username or password does not match"), domainErrors.NotAuthorized)
	}
	s.Logger.Info("User login attempt", zap.String("username", username))

	// Single Sign On offline user
	if os.Getenv("SERVER_SINGLE_SIGN_ON") == "true" {
		s.sessionManager.NotifyOtherDevicesOffline(user.ID, ginCtx.Query("deviceId"))
	}

	isAuthenticated := sharedUtil.CheckPasswordHash(password, user.HashPassword)
	if !isAuthenticated {
		s.Logger.Warn("Login failed: invalid password", zap.String("username", username))
		return nil, nil, nil, domainErrors.NewAppError(errors.New("username or password does not match"), domainErrors.NotAuthorized)
	}

	var role domainRole.Role
	var roleId int64
	if len(user.Roles) > 0 {
		roleId = user.Roles[0].ID
		role = user.Roles[0]
	}
	accessTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "access")
	if err != nil {
		s.Logger.Error("Error generating access token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, nil, err
	}
	refreshTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "refresh")
	if err != nil {
		s.Logger.Error("Error generating refresh token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, nil, err
	}

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		RefreshToken:              refreshTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		ExpirationRefreshDateTime: refreshTokenClaims.ExpirationTime,
	}
	ctx := context.Background()
	s.RedisClient.Set(ctx, GetUserTokenKey(user.ID), accessTokenClaims.Token, UserTokenExpireDuration)
	s.RedisClient.Set(ctx, GetUserRefreshTokenKey(user.ID), refreshTokenClaims.Token, RefreshTokenExpireDuration)

	s.Logger.Info("User login successful", zap.String("username", username), zap.Int64("userID", user.ID))
	return user, authTokens, &role, nil
}
func (s *AuthUseCase) AccessTokenByRefreshToken(refreshToken string) (*domainUser.User, *AuthTokens, error) {
	s.Logger.Info("Refreshing access token")

	// 首先验证刷新令牌本身的有效性
	claimsMap, err := s.JWTService.GetClaimsAndVerifyToken(refreshToken, "refresh")
	if err != nil {
		s.Logger.Error("Error verifying refresh token", zap.Error(err))
		return nil, nil, err
	}

	userID := int(claimsMap["id"].(float64))
	roleId := int64(claimsMap["role_id"].(float64))

	// 检查刷新令牌是否在黑名单中
	exists, err := s.jwtBlacklistRepository.IsJwtInBlacklist(refreshToken)
	if err != nil {
		s.Logger.Error("Error checking refresh token in blacklist", zap.Error(err))
		return nil, nil, domainErrors.NewAppError(err, domainErrors.TokenError)
	}

	if exists {
		s.Logger.Warn("Refresh token is in blacklist", zap.Int("userID", userID))
		return nil, nil, domainErrors.NewAppError(errors.New("refresh token has been revoked"), domainErrors.TokenError)
	}

	// 检查是否启用了单点登录，并验证是否为当前活跃会话
	if os.Getenv("SERVER_SINGLE_SIGN_ON") == "true" {
		// 获取用户当前的活跃刷新令牌
		currentRefreshToken, err := s.RedisClient.Get(context.Background(), GetUserRefreshTokenKey(int64(userID))).Result()
		if err == nil && currentRefreshToken != refreshToken {
			s.Logger.Warn("Refresh token has been replaced", zap.Int("userID", userID))
			// 将旧的刷新令牌加入黑名单
			_ = s.jwtBlacklistRepository.AddToBlacklist(refreshToken)
			return nil, nil, domainErrors.NewAppError(errors.New("refresh token has been replaced"), domainErrors.TokenError)
		}
	}

	// 获取用户信息
	user, err := s.UserRepository.GetByID(userID)
	if err != nil {
		s.Logger.Error("Error getting user for token refresh", zap.Error(err), zap.Int("userID", userID))
		return nil, nil, err
	}

	// 生成新的访问令牌
	accessTokenClaims, err := s.JWTService.GenerateJWTToken(user.ID, roleId, "access")
	if err != nil {
		s.Logger.Error("Error generating new access token", zap.Error(err), zap.Int64("userID", user.ID))
		return nil, nil, err
	}

	var expTime = int64(claimsMap["exp"].(float64))

	authTokens := &AuthTokens{
		AccessToken:               accessTokenClaims.Token,
		ExpirationAccessDateTime:  accessTokenClaims.ExpirationTime,
		RefreshToken:              refreshToken,
		ExpirationRefreshDateTime: time.Unix(expTime, 0),
	}
	ctx := context.Background()
	s.RedisClient.Set(ctx, GetUserTokenKey(user.ID), accessTokenClaims.Token, UserTokenExpireDuration)

	s.Logger.Info("Access token refreshed successfully", zap.Int64("userID", user.ID))
	return user, authTokens, nil
}

// Register implements IAuthUseCase.
func (s *AuthUseCase) Register(user RegisterUser) (*domain.CommonResponse[SecurityRegisterUser], error) {
	// user is exist
	whereCondition := make(map[string]interface{}, 3)
	whereCondition["user_name"] = user.UserName
	dbUser, err := s.UserRepository.GetOneByMap(whereCondition)
	if err != nil {
		return nil, err
	}
	if dbUser.ID != 0 {
		return nil,
			domainErrors.NewAppError(errors.New("The user already exists"), domainErrors.UserExists)
	}
	userRepo := domainUser.User{
		UserName: user.UserName,
		Email:    user.Email,
		Password: user.Password,
	}
	// password to has
	hash, err := sharedUtil.StringToHash(user.Password)
	if err != nil {
		return &domain.CommonResponse[SecurityRegisterUser]{}, err
	}
	userRepo.HashPassword = string(hash)

	// generate uuid
	userRepo.UUID = uuid.New().String()
	userRepo.Status = 1

	res, err := s.UserRepository.Create(&userRepo)
	if err != nil {
		return &domain.CommonResponse[SecurityRegisterUser]{}, err
	}

	return &domain.CommonResponse[SecurityRegisterUser]{
		Data: SecurityRegisterUser{
			Data: DataUserAuthenticated{
				ID:       res.ID,
				UUID:     res.UUID,
				UserName: res.UserName,
				NickName: res.NickName,
				Email:    res.Email,
				Status:   res.Status,
			},
		},
		Status:  0,
		Message: "success",
	}, nil

}

func (s *AuthUseCase) Logout(jwtToken string) (*domain.CommonResponse[string], error) {
	// 解析token获取用户ID
	claimsMap, err := s.JWTService.GetClaimsAndVerifyToken(jwtToken, "access")
	if err != nil {
		// 如果是访问令牌解析失败，尝试解析刷新令牌
		claimsMap, err = s.JWTService.GetClaimsAndVerifyToken(jwtToken, "refresh")
		if err != nil {
			return nil, domainErrors.NewAppError(err, domainErrors.TokenError)
		}
	}

	userID := int64(claimsMap["id"].(float64))

	// 检查token是否已经在黑名单中
	exist, err := s.jwtBlacklistRepository.IsJwtInBlacklist(jwtToken)
	if err != nil {
		return nil, domainErrors.NewAppError(err, domainErrors.TokenError)
	}
	if exist {
		return nil, domainErrors.NewAppError(errors.New("the user logout already"), domainErrors.TokenError)
	}
	// 将对应的访问令牌和刷新令牌都加入黑名单
	ctx := context.Background()
	// 将当前令牌加入黑名单
	err = s.jwtBlacklistRepository.AddToBlacklist(jwtToken)
	if err != nil {
		return nil, domainErrors.NewAppError(err, domainErrors.TokenError)
	}

	// 获取并加入刷新令牌到黑名单
	if refreshToken, err := s.RedisClient.Get(ctx, GetUserRefreshTokenKey(userID)).Result(); err == nil {
		if refreshToken != "" && refreshToken != jwtToken {
			_ = s.jwtBlacklistRepository.AddToBlacklist(refreshToken)
		}
	}

	// 如果启用了单点登录，清除Redis中的活跃令牌记录
	if os.Getenv("SERVER_SINGLE_SIGN_ON") == "true" {
		ctx := context.Background()
		s.RedisClient.Del(ctx, GetUserTokenKey(userID))
		s.RedisClient.Del(ctx, GetUserRefreshTokenKey(userID))
	}

	return &domain.CommonResponse[string]{Data: "true", Status: 0, Message: "success"}, nil
}
