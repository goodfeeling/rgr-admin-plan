package middlewares

import (
	"context"
	"net/http"
	"os"
	"strings"

	authUseCase "github.com/gbrayhan/microservices-go/src/application/services/auth"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// AuthJWTMiddlewareWithRedis 使用Redis的认证中间件
func AuthJWTMiddlewareWithRedis(redisClient *redis.Client, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided"})
			c.Abort()
			return
		}

		if !CommonVerifyWithRedis(c, tokenString, redisClient, db) {
			return
		}

		c.Next()
	}
}

// OptionalAuthMiddlewareWithRedis 可选认证中间件
func OptionalAuthMiddlewareWithRedis(redisClient *redis.Client, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.Next()
			return
		}

		CommonVerifyWithRedis(c, tokenString, redisClient, db)
		c.Next()
	}
}

// UrlAuthMiddlewareWithRedis URL参数认证中间件
func UrlAuthMiddlewareWithRedis(redisClient *redis.Client, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.Query("token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token not provided1"})
			c.Abort()
			return
		}
		if !CommonVerifyWithRedis(c, tokenString, redisClient, db) {
			return
		}

		c.Next()
	}
}

// CommonVerifyWithRedis 带Redis验证的通用验证函数
func CommonVerifyWithRedis(c *gin.Context, tokenString string, redisClient *redis.Client, db *gorm.DB) bool {
	accessSecret := os.Getenv("JWT_ACCESS_SECRET")
	if accessSecret == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "JWT_ACCESS_SECRET not configured"})
		c.Abort()
		return false
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return []byte(accessSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		c.Abort()
		return false
	}

	// check token if in blacklist
	exists, err := IsTokenInBlacklist(db, tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		c.Abort()
		return false
	}

	if exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
		c.Abort()
		return false
	}

	// 检查用户当前有效token（实现单点登录）
	if userID, ok := claims["id"].(float64); ok && os.Getenv("SERVER_SINGLE_SIGN_ON") == "true" {
		currentToken, err := redisClient.Get(context.Background(), authUseCase.GetUserTokenKey(int64(userID))).Result()
		if err == nil && currentToken != tokenString {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been replaced"})
			c.Abort()
			return false
		}
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < jwt.TimeFunc().Unix() {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return false
		}
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		c.Abort()
		return false
	}

	if t, ok := claims["type"].(string); ok {
		if t != "access" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Token type mismatch"})
			c.Abort()
			return false
		}
	} else {
		c.JSON(http.StatusForbidden, gin.H{"error": "Missing token type"})
		c.Abort()
		return false
	}

	if idFloat, ok := claims["id"].(float64); ok {
		id := int(idFloat)
		if id == 0 {
			c.JSON(http.StatusForbidden, gin.H{"error": "Missing or invalid user id"})
			c.Abort()
			return false
		}
		c.Set("user_id", id)
	}

	if roleIdFloat64, ok := claims["role_id"].(float64); ok {
		id := int64(roleIdFloat64)
		if id == 0 {
			id = -1
		}
		c.Set("role_id", id)
	}

	return true
}
