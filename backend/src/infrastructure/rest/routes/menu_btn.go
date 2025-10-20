package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func MenuBtnRouters(
	router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.MenuBtnModule.Controller
	u := router.Group("/menu_btn")
	u.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	u.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		u.POST("", controller.NewMenuBtn)
		u.GET("", controller.GetAllMenuBtns)
		u.GET("/:id", controller.GetMenuBtnsByID)
		u.PUT("/:id", controller.UpdateMenuBtn)
		u.DELETE("/:id", controller.DeleteMenuBtn)
	}
}
