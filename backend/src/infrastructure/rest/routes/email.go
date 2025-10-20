package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
)

func EmailRouters(router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.EmailModule.Controller
	u := router.Group("/email")
	{
		u.POST("/forget-password", controller.SendForgetPasswordEmail)
	}
}
