package controllers

import (
	"authentication/config"
	"authentication/helpers"
	models "authentication/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
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

		if userClaims.Role != "ADMIN" && userClaims.Role != "SUPER_ADMIN" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}

		// ✅ Use this query
		var users []models.User
		err := config.DB.Select(&users, `
			SELECT id, name, email, phone, address, role, status, subscription_status, password, created_at, updated_at, deleted_at
			FROM users
		`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}
