package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuGroupRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.MenuGroupModule.Controller
	u := router.Group("/menu_group")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.POST("", controller.NewMenuGroup)
		u.GET("", controller.GetAllMenuGroups)
		u.GET("/:id", controller.GetMenuGroupsByID)
		u.PUT("/:id", controller.UpdateMenuGroup)
		u.DELETE("/:id", controller.DeleteMenuGroup)
		u.GET("/search", controller.SearchPaginated)
		u.GET("/search-property", controller.SearchByProperty)
	}
}
