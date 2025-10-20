package routes

import (
	"net/http"

	"github.com/gbrayhan/microservices-go/src/infrastructure/di"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func ApplicationRouter(router *gin.Engine, appContext *di.ApplicationContext) {
	v1 := router.Group("/v1")

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Service is running",
		})
	})

	AuthRoutes(v1, appContext)
	UserRoutes(v1, appContext)
	UploadRoutes(v1, appContext)
	RoleRoutes(v1, appContext)
	ApiRouters(v1, router, appContext)
	OperationRouters(v1, appContext)
	DictionaryRouters(v1, appContext)
	DictionaryDetailRouters(v1, appContext)
	MenuRouters(v1, appContext)
	MenuGroupRouters(v1, appContext)
	MenuBtnRouters(v1, appContext)
	MenuParameterRouters(v1, appContext)
	FileRouters(v1, appContext)

	ScheduledTaskRouters(v1, appContext)
	ConfigRouters(v1, appContext)
	TaskExecutionLogRouters(v1, appContext)
	EmailRouters(v1, appContext)
	CaptchaRoutes(v1, appContext)
}
