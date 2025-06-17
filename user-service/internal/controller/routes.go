package controller

import (
	"github.com/gin-gonic/gin"
	"user-service/internal/controller/user"
	"user-service/internal/service"
)

func RouteAPI(route *gin.Engine, services *service.Service) {
	userContr := user.NewUserController(services.UserService)
	//docs.SwaggerInfo.BasePath = "/api"

	route.OPTIONS("/*path", func(c *gin.Context) {
		c.AbortWithStatus(204)
	})

	route.Use(gin.Logger())
	//route.GET("/swagger", ginSwagger.WrapHandler(swaggerFiles.Handler))
	route.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := route.Group("/api")
	//api.Use(LoggerMiddleware())

	userApi := api.Group("/user")
	userApi.Use(AuthMiddleware())
	userApi.GET("/get", userContr.Get)
	userApi.GET("/validate", userContr.Validate)
	userApi.POST("/get/like", userContr.GetUsersByUsername)
	userApi.POST("/update", userContr.UpdateUserByID)

	outsideApi := api.Group("/outside/user")

	outsideApi.GET("/get/:id", userContr.OutsideGetUserByID)
}
