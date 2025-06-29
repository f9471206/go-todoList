// routes/admin.go
package routes

import (
	"todolist/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(r *gin.RouterGroup) {

	controller := controllers.AuthController{}
	admin := r.Group("/")
	{
		admin.POST("/login", controller.Login)
		admin.POST("/register", controller.Register)
	}
}
