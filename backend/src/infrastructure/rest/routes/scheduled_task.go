package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func ScheduledTaskRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.ScheduledTaskModule.Controller
	u := router.Group("/scheduled_task")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.POST("", controller.NewScheduledTask)
		u.GET("", controller.GetAllScheduledTasks)
		u.GET("/:id", controller.GetScheduledTaskByID)
		u.PUT("/:id", controller.UpdateScheduledTask)
		u.DELETE("/:id", controller.DeleteScheduledTask)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
		u.POST("/delete-batch", controller.DeleteScheduledTasks)
		u.POST("/enable/:id", controller.EnableTaskById)
		u.POST("/disable/:id", controller.DisableTaskById)
		u.POST("/reload", controller.ReloadAllTasks)
	}
}
