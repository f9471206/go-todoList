// routes/admin.go
package routes

import (
	"todolist/controllers"
	"todolist/middleware"

	"github.com/gin-gonic/gin"
)

func MemberRoutes(r *gin.RouterGroup) {

	controller := controllers.MenberController{}
	member := r.Group("/member", middleware.JwtAuthMiddleware())

	member.GET("/", middleware.RequireRoles("Admin"), controller.Index)
	member.GET("/:id", middleware.RequireRoles("Admin"), controller.Show)
	member.PUT("/:id", middleware.RequireRoles("Admin"), controller.Edit)

}
