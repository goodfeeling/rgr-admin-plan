package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func RoleRoutes(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.RoleModule.Controller
	u := router.Group("/role")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.POST("", controller.NewRole)
		u.GET("", controller.GetAllRoles)
		u.GET("/:id", controller.GetRolesByID)
		u.PUT("/:id", controller.UpdateRole)
		u.DELETE("/:id", controller.DeleteRole)
		u.GET("/tree", controller.GetTreeRoles)
		u.GET("/:id/setting", controller.GetRoleSetting)
		u.POST("/:id/menu", controller.UpdateRoleMenuIds)
		u.POST("/:id/api", controller.BindApiRule)
		u.POST("/:id/menu-btns", controller.BindRoleMenuBtns)
	}
}
