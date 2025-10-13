package routes

import (
	"authentication/controllers"
	"authentication/helpers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api")
	{
		api.POST("/signup", controllers.Signup())
		api.POST("/login", controllers.Login())

		userRoutes := api.Group("/users")
		userRoutes.Use(helpers.AuthMiddleware())
		{
			userRoutes.GET("/", helpers.AuthorizeRole("ADMIN", "SUPER_ADMIN"), controllers.GetUsers())
			userRoutes.GET("/:id", controllers.GetUser())
		}
	}
}
