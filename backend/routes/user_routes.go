package routes

import (
	"taskflow/controllers"
	"taskflow/middleware"

	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	protected := r.Group("/")
	protected.Use(middleware.AuthMiddleware())
	protected.GET("/users", controllers.ListUsers)
}
