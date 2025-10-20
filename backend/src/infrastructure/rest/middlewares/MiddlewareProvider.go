package middlewares

import (
	"errors"

	"github.com/gbrayhan/microservices-go/src/infrastructure/repository/psql/jwt_blacklist"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type MiddlewareProvider struct {
	RedisClient *redis.Client
	DB          *gorm.DB
}

func NewMiddlewareProvider(redisClient *redis.Client, db *gorm.DB) *MiddlewareProvider {
	return &MiddlewareProvider{
		RedisClient: redisClient,
		DB:          db,
	}
}

func (mp *MiddlewareProvider) AuthJWTMiddleware() gin.HandlerFunc {
	return AuthJWTMiddlewareWithRedis(mp.RedisClient, mp.DB)
}

func (mp *MiddlewareProvider) OptionalAuthMiddleware() gin.HandlerFunc {
	return OptionalAuthMiddlewareWithRedis(mp.RedisClient, mp.DB)
}

func (mp *MiddlewareProvider) UrlAuthMiddleware() gin.HandlerFunc {
	return UrlAuthMiddlewareWithRedis(mp.RedisClient, mp.DB)
}

func (mp *MiddlewareProvider) AuthResetPasswordMiddleware() gin.HandlerFunc {
	return AuthResetPassword(mp.RedisClient, mp.DB)
}

// IsTokenInBlacklist 检查token是否在黑名单中
func IsTokenInBlacklist(db *gorm.DB, tokenString string) (bool, error) {
	var blacklist jwt_blacklist.JwtBlacklist
	err := db.Where("jwt = ? AND deleted_at IS NULL", tokenString).First(&blacklist).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
