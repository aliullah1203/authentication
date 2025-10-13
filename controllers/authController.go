package controllers

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/models"
	"net/http"
	"time"

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

		if err := helpers.ValidateStruct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Check duplicate
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

		user.ID = uuid.New()
		user.Role = "CUSTOMER"
		user.Status = "ACTIVE"
		user.SubscriptionStatus = "SUBSCRIBED"
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		user.Password = helpers.HashPassword(user.Password)

		_, err = config.DB.NamedExec(`INSERT INTO users (id, name, email, phone, address, role, status, subscription_status, password, created_at, updated_at) 
			VALUES (:id,:name,:email,:phone,:address,:role,:status,:subscription_status,:password,:created_at,:updated_at)`, &user)
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

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.User
		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var user models.User
		err := config.DB.Get(&user, "SELECT * FROM users WHERE email=$1", req.Email)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		if !helpers.VerifyPassword(user.Password, req.Password) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		accessToken, refreshToken, err := helpers.GenerateTokens(user.ID.String(), user.Role)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"user":          user,
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})
	}
}

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

		var user models.User
		err := config.DB.Get(&user, "SELECT * FROM users WHERE id=$1", requestedUserId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"user": user})
	}
}

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

		var users []models.User
		err := config.DB.Select(&users, "SELECT * FROM users")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"users": users})
	}
}
