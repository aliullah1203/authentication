package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"authentication/config"
	"authentication/internal/helpers"
	models "authentication/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Signup(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, `{"error":"invalid JSON: `+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	// Check duplicate email or phone
	var count int
	err := config.DB.Get(&count, "SELECT COUNT(*) FROM users WHERE email=$1 OR phone=$2", user.Email, user.Phone)
	if err != nil {
		http.Error(w, `{"error":"database error: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}
	if count > 0 {
		http.Error(w, `{"error":"email or phone already exists"}`, http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"failed to hash password"}`, http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Fill other fields
	user.ID = uuid.New()
	user.Role = "CUSTOMER"
	user.Status = "ACTIVE"
	user.SubscriptionStatus = "SUBSCRIBED"
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Insert user into DB
	_, err = config.DB.NamedExec(`INSERT INTO users 
		(id, name, email, phone, address, role, status, subscription_status, password, created_at, updated_at) 
		VALUES 
		(:id,:name,:email,:phone,:address,:role,:status,:subscription_status,:password,:created_at,:updated_at)`, &user)
	if err != nil {
		http.Error(w, `{"error":"insert error: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	token, err := helpers.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		http.Error(w, `{"error":"token generation failed"}`, http.StatusInternalServerError)
		return
	}

	// Send JSON response
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User created successfully",
		"token":   token,
	})
}
