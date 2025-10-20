package middlewares

import (
	"fmt"
	"net/http"

	"github.com/casbin/casbin/v2"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers"
	"github.com/gin-gonic/gin"
)

// CasbinMiddleware 创建一个Casbin权限验证中间件
func CasbinMiddleware(enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取应用上下文
		appCtx := controllers.NewAppUtils(c)
		// 获取角色ID
		roleId, ok := appCtx.GetRoleID()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - Role not found",
			})
			return
		}

		userId, ok := appCtx.GetUserID()
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Unauthorized - User not found",
			})
			return
		}
		// super user jump verify
		if userId == 1 {
			c.Next()
			return
		}

		// 将用户ID作为角色进行权限验证
		role := fmt.Sprintf("%d", roleId) // v0字段存储的是角色ID
		path := c.Request.URL.Path        // v1字段存储的是API路径
		method := c.Request.Method        // v3字段存储的是请求方法

		// 使用Casbin进行权限检查
		ok, err := enforcer.Enforce(role, path, method)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error during authorization",
			})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Forbidden - Insufficient permissions",
			})
			return
		}

		// 权限验证通过，继续处理请求
		c.Next()
	}
}
