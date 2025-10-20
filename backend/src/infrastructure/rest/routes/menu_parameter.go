package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuParameterRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.MenuParameterModule.Controller
	u := router.Group("/menu_parameter")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.POST("", controller.NewMenuParameter)
		u.GET("", controller.GetAllMenuParameters)
		u.GET("/:id", controller.GetMenuParametersByID)
		u.PUT("/:id", controller.UpdateMenuParameter)
		u.DELETE("/:id", controller.DeleteMenuParameter)
	}
}
