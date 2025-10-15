package controllers

import (
	"authentication/config"
	models "authentication/user"
	"net/http"
	"time"

	"authentication/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check duplicate email or phone
		var count int
		err := config.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1 OR phone=$2", user.Email, user.Phone)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email or phone already exists"})
			return
		}

		// Fill other fields
		user.ID = uuid.New()
		user.Role = "CUSTOMER"
		user.Status = "ACTIVE"
		user.SubscriptionStatus = "SUBSCRIBED"
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()

		// Store password as plain text
		user.Password = user.Password

		_, err = config.DB.NamedExec(`INSERT INTO users 
		(id, name, email, phone, address, role, status, subscription_status, password, created_at, updated_at) 
		VALUES 
		(:id,:name,:email,:phone,:address,:role,:status,:subscription_status,:password,:created_at,:updated_at)`, &user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		accessToken, refreshToken, err := helpers.GenerateTokens(user.ID.String(), user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":       "User created successfully",
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}
