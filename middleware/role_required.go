package middleware

import (
	"net/http"
	"strings"
	"todolist/config"
	"todolist/models"
	"todolist/utils"

	"github.com/gin-gonic/gin"
)

// middlewares/role_required.go
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {

		userID := c.Request.Context().Value(utils.UserIDKey)

		if userID == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		var user models.User
		if err := config.DB.Preload("Roles").First(&user, userID).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "user not found"})
			return
		}

		for _, role := range user.Roles {
			for _, allowed := range allowedRoles {
				if strings.EqualFold(role.Name, allowed) {
					c.Next()
					return
				}
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
	}
}
