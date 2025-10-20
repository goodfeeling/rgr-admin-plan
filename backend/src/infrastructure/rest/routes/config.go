package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func ConfigRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.ConfigModule.Controller
	u := router.Group("/config")
	u.GET("/site", controller.GetConfigBySite)

	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.GET("", controller.GetAllConfigs)
		u.PUT("/:module", controller.UpdateConfig)
		u.GET("/:module", controller.GetConfigByModule)
	}
}
