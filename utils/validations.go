package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BindAndValidate(c *gin.Context, obj interface{}) bool {
	if err := c.ShouldBindJSON(obj); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return false
	}
	return true
}
