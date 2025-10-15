package controllers

import (
	"authentication/config"
	models "authentication/user"
	"net/http"

	"authentication/helpers"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	var req models.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	var user models.User
	// Only select required columns
	query := `SELECT id, name, email, role, password FROM users WHERE email=$1 LIMIT 1`
	err := config.DB.Get(&user, query, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare plain passwords
	if user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

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
