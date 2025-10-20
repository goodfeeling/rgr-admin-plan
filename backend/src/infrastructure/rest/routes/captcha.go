package routes

import (
	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
)

func CaptchaRoutes(router *gin.RouterGroup, appContext *di.ApplicationContext) {
	controller := appContext.CaptchaModule.Controller
	r := router.Group("/captcha")
	{
		r.GET("/generate", controller.Generate)
	}

}
