package middleware

import (
	"context"
	"net/http"
	"strings"
	"todolist/config"
	"todolist/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// JwtAuthMiddleware 驗證 JWT Token 的中間件
func JwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Authorization header 取出 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 驗證簽名算法
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			// 回傳簽名用的密鑰（跟簽發 token 時用的一樣）
			return []byte(config.JWTSecret), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["user_id"]

			utils.Logger.Info("JWT驗證成功",
				zap.Any("user_id", userID),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("client_ip", c.ClientIP()),
			)

			ctx := context.WithValue(c.Request.Context(), utils.UserIDKey, userID)
			c.Request = c.Request.WithContext(ctx)
			c.Next()
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
	}
}
