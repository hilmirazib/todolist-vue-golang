package routes

import (
	"github.com/gin-gonic/gin"
	"my-todolist/handlers"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.GET("/todos", handlers.ListTodos)
		api.POST("/todos", handlers.CreateTodo)
		api.PATCH("/todos/:id", handlers.UpdateTodo)
		api.DELETE("/todos/:id", handlers.DeleteTodo)
	}
}
