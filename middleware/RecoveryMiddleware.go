package middleware

import (
	"net/http"
	"todolist/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				userID := c.Request.Context().Value(utils.UserIDKey)

				utils.Logger.Error("panic recovered",
					zap.Any("error", r),
					zap.Any("user_id", userID),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
				)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}
