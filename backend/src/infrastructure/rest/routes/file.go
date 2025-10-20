package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func FileRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.FileModule.Controller
	u := router.Group("/file")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.POST("", controller.NewFile)
		u.GET("", controller.GetAllFiles)
		u.GET("/:id", controller.GetFilesByID)
		u.PUT("/:id", controller.UpdateFile)
		u.DELETE("/:id", controller.DeleteFile)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
	}
}
