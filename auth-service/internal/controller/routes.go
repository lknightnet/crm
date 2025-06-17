package controller

import (
	"auth-service/internal/controller/auth"
	"auth-service/internal/service"
	"github.com/gin-gonic/gin"
)

func RouteAPI(route *gin.Engine, services *service.Service) {
	authContr := auth.NewAuthController(services.AuthService)

	//docs.SwaggerInfo.BasePath = "/api"

	route.Use(gin.Logger())
	//route.GET("/swagger", ginSwagger.WrapHandler(swaggerFiles.Handler))
	route.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := route.Group("/api")
	//api.Use(LoggerMiddleware())

	//updateApi := api.Group("/update")
	//
	//updateApi.POST("/name", nil)

	authApi := api.Group("/auth")

	authApi.POST("/login", authContr.Login)
	authApi.POST("/signup", authContr.Signup)
	authApi.GET("/verify", authContr.Verify)
	authApi.POST("/refresh", authContr.Refresh)
}
