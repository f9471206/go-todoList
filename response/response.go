package response

import (
	"net/http"
	"todolist/utils"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func SuccessWithPagination[T any](c *gin.Context, data *utils.PaginatedResult[T]) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

func SuccessWithMessage(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

func Error(c *gin.Context, statusCode int, msg string) {
	c.JSON(statusCode, gin.H{
		"success": false,
		"error":   msg,
	})
}
