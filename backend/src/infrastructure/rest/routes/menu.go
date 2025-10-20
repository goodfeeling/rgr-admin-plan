package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.MenuModule.Controller
	middlewareProvider := appContext.MiddlewareProvider
	u := router.Group("/menu")
	u.GET("/user", middlewareProvider.OptionalAuthMiddleware(), controller.GetUserMenus)

	protected := u.Group("")
	protected.Use(middlewareProvider.AuthJWTMiddleware())
	protected.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		protected.POST("", controller.NewMenu)
		protected.GET("", controller.GetAllMenus)
		protected.GET("/:id", controller.GetMenusByID)
		protected.PUT("/:id", controller.UpdateMenu)
		protected.DELETE("/:id", controller.DeleteMenu)
	}

}
