package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func OperationRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.OperationModule.Controller
	u := router.Group("/operation")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.GET("", controller.GetAllOperations)
		u.GET("/:id", controller.GetOperationsByID)
		u.DELETE("/:id", controller.DeleteOperation)
		u.POST("/delete-batch", controller.DeleteOperations)
		u.GET("/search", controller.SearchPaginated)
	}
}
