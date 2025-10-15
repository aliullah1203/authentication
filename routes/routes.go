package routes

import (
	"authentication/controllers"
	"authentication/helpers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		// Public routes
		api.POST("/signup", controllers.Signup())
		api.POST("/login", controllers.Login)

		// Ping
		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})

		// Protected routes
		userRoutes := api.Group("/users")
		userRoutes.Use(helpers.AuthMiddleware())
		{
			// only ADMIN or SUPER_ADMIN allowed to list users
			userRoutes.GET("/", helpers.AuthorizeRole("ADMIN", "SUPER_ADMIN"), controllers.GetUsers())
			userRoutes.GET("/:id", controllers.GetUser())
		}
	}
}
