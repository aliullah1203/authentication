package controllers

import (
	"authentication/config"
	"authentication/helpers"
	models "authentication/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedUserId := c.Param("id")

		claims, exists := c.Get("claims")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		userClaims, ok := claims.(*helpers.Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		if userClaims.Role != "ADMIN" && userClaims.UserID != requestedUserId {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// ✅ Use this query
		var user models.User
		err := config.DB.Get(&user, `
			SELECT id, name, email, phone, address, role, status, subscription_status, password, created_at, updated_at, deleted_at
			FROM users WHERE id=$1
		`, requestedUserId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}
