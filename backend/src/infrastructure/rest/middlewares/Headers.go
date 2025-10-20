package middlewares

import (
	"strings"
	"time"

	sharedUtil "github.com/gbrayhan/microservices-go/src/shared/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// security headers
func SecurityHeaders() gin.HandlerFunc {

	return func(c *gin.Context) {
		c.Header("X-Frame-Options", "SAMEORIGIN")
		c.Header("Cache-Control", "no-cache, no-store")
		c.Header("Pragma", "no-cache")
		c.Header("Expires", "0")

		c.Next()
	}

}

// cors header set
func CorsHeader() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(sharedUtil.GetEnv("CORS_ALLOWED_ORIGINS", ""), ","),
		AllowMethods:     strings.Split(sharedUtil.GetEnv("CORS_ALLOWED_METHODS", ""), ","),
		AllowHeaders:     strings.Split(sharedUtil.GetEnv("CORS_ALLOWED_HEADERS", ""), ","),
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})
}
