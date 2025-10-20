package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func TaskExecutionLogRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.TaskExecutionLogModule.Controller
	u := router.Group("/task_execution_log")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.GET("/search", controller.SearchPaginated)
	}
}
