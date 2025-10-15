package controllers

import (
	"authentication/config"
	"authentication/helpers"
	models "authentication/user"
	"net/http"

	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req models.User

	// Parse JSON body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Fetch user from DB
	var user models.User
	query := "SELECT * FROM users WHERE email=$1 LIMIT 1"
	err := config.DB.Get(&user, query, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Verify password (assuming bcrypt)
	if !helpers.VerifyPassword(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Generate JWT tokens
	accessToken, refreshToken, err := helpers.GenerateTokens(user.ID.String(), user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
