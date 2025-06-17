package user

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func (userc *UserController) AuthMiddleware() gin.HandlerFunc {
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
