package controller

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"project-service/pkg/tg"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(401, gin.H{
				"error": "Authorization token missing or malformed",
			})
			c.Abort()
			return
		}

		token := authHeader[len("Bearer "):]

		if token == "" {
			c.JSON(401, gin.H{
				"error": "Authorization token is empty",
			})
			c.Abort()
			return
		}

		c.Set("token", token)

		c.Next()
	}
}

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestInfo := make(map[string]interface{})

		// Сохраняем Query параметры
		requestInfo["query"] = c.Request.URL.Query()

		// Сохраняем Headers
		requestInfo["headers"] = c.Request.Header

		// Сохраняем Method и Path
		requestInfo["method"] = c.Request.Method

		// Сохраняем IP
		requestInfo["client_ip"] = c.ClientIP()

		// Сохраняем Body (осторожно: body можно читать только один раз!)
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Восстановим Body, чтобы дальше работало

			requestInfo["body"] = string(bodyBytes)
		}

		tg.SendInfo(requestInfo, c.Request.URL.Path)

		c.Next()
	}
}
