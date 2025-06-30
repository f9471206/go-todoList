package routes

import (
	"todolist/controllers"
	"todolist/middleware"

	"github.com/gin-gonic/gin"
)

func TodoRoutes(r *gin.RouterGroup) {
	todoTypeController := controllers.TodoTypeController{}
	todoListController := controllers.TodoListController{}
	todoListDetailsController := controllers.TodoListDetailsController{}

	todo := r.Group("/todo", middleware.JwtAuthMiddleware())
	{
		todo.POST("/type", todoTypeController.Create)
		todo.GET("/type", todoTypeController.Index)
		todo.GET("/type/:id", todoTypeController.Show)
		todo.PUT("/type/:id", todoTypeController.Edit)
		todo.DELETE("/type/:id", todoTypeController.Delete)

		todo.POST("/list", todoListController.Create)
		todo.GET("/list", todoListController.Index)
		todo.GET("/list/:id", todoListController.Show)
		todo.PUT("/list/:id", todoListController.Edit)
		todo.DELETE("/list/:id", todoListController.Delete)

		todo.POST("/list/details", todoListDetailsController.Create)
		todo.PUT("/list/details/:id", todoListDetailsController.Edit)
		todo.DELETE("list/details/:id", todoListDetailsController.Delete)
	}
}
