package controller

import (
	"auth-service/pkg/tg"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
)

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
