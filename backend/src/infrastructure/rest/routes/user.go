package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gbrayhan/microservices-go/src/infrastructure/rest/middlewares"
	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.UserModule.Controller
	middlewareProvider := appContext.MiddlewareProvider

	u := router.Group("/user")

	u.POST("/change-password", middlewareProvider.AuthResetPasswordMiddleware(), controller.ChangePassword)

	protected := u.Group("")
	protected.Use(middlewareProvider.AuthJWTMiddleware())
	protected.Use(middlewares.CasbinMiddleware(appContext.Enforcer))
	{
		protected.POST("", controller.NewUser)
		protected.GET("", controller.GetAllUsers)
		protected.GET("/:id", controller.GetUsersByID)
		protected.PUT("/:id", controller.UpdateUser)
		protected.DELETE("/:id", controller.DeleteUser)
		protected.GET("/search", controller.SearchPaginated)
		protected.GET("/search-property", controller.SearchByProperty)
		protected.POST("/:id/role", controller.UserBindRoles)
		protected.POST("/:id/reset-password", controller.ResetPassword)
		protected.POST("/:id/edit-password", controller.EditPassword)
	}

}
