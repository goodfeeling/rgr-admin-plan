package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func UploadRoutes(router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.UploadModule.Controller
	u := router.Group("/upload")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.POST("/single", controller.Single)
		u.POST("/multiple", controller.Multiple)
		u.GET("/sts-token", controller.GetSTSToken)
		u.GET("/refresh-sts", controller.RefreshSTSToken)
	}
}
