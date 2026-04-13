package routes

import (
	"taskflow/controllers"
	"taskflow/middleware"

	"github.com/gin-gonic/gin"
)

func ProjectRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())

	protected.POST("/projects", controllers.CreateProject)
	protected.GET("/projects", controllers.ListProjects)
	protected.GET("/projects/:id", controllers.GetProject)
	protected.PATCH("/projects/:id", controllers.UpdateProject)
	protected.DELETE("/projects/:id", controllers.DeleteProject)
	protected.GET("/projects/:id/stats", controllers.ProjectStats)

	protected.POST("/projects/:id/tasks", controllers.CreateTask)
	protected.GET("/projects/:id/tasks", controllers.ListTasks)

	protected.PATCH("/tasks/:id", controllers.UpdateTask)
	protected.DELETE("/tasks/:id", controllers.DeleteTask)
}
