package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/controllers/api"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func ApiRouters(
	router *gin.RouterGroup,
	routerEngine *gin.Engine, appContext *di.ApplicationContext) {
	controller := appContext.ApiModule.Controller
	// 用户获取接口列表
	if routerSetter, ok := controller.(api.RouterSetter); ok {
		routerSetter.SetRouter(routerEngine)
	}
	u := router.Group("/api")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.POST("", controller.NewApi)
		u.GET("", controller.GetAllApis)
		u.GET("/:id", controller.GetApisByID)
		u.PUT("/:id", controller.UpdateApi)
		u.DELETE("/:id", controller.DeleteApi)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
		u.POST("/delete-batch", controller.DeleteApis)
		u.GET("/group-list", controller.GetApisGroup)
		u.POST("/synchronize", controller.SynchronizeRouterToApi)
		u.GET("/excel/template", controller.DownloadTemplate)
		u.POST("/excel/import", controller.Import)
		u.GET("/excel/export", controller.Export)
	}
}
