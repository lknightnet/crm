package controller

import (
	"github.com/gin-gonic/gin"
	"jwt-service/internal/service"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"jwt-service/docs"
)

func RouteAPI(route *gin.Engine, services *service.Service) {
	jwtContr := newJWTController(services.JWTService)

	docs.SwaggerInfo.BasePath = "/api"

	route.Use(gin.Logger())
	route.GET("/swagger", ginSwagger.WrapHandler(swaggerFiles.Handler))
	route.GET("/healthcheck", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := route.Group("/api")
	//api.Use(LoggerMiddleware())
	api.POST("/generate", jwtContr.Generate)
	api.POST("/validate", jwtContr.Validate)
	api.POST("/refresh", jwtContr.Refresh)
}
