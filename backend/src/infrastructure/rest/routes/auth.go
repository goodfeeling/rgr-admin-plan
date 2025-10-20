package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.AuthModule.Controller
	routerAuth := router.Group("/auth")
	{
		routerAuth.POST("/signin", controller.Login)
		routerAuth.POST("/signup", controller.Register)
		routerAuth.POST("/access-token", controller.GetAccessTokenByRefreshToken)
	}
	loginAuth := routerAuth.Use(appContext.MiddlewareProvider.AuthJWTMiddleware())
	{
		loginAuth.POST("/switch-role", controller.SwitchRole)
		loginAuth.GET("/logout", controller.Logout)
	}
}
