package controllers

import (
	"authentication/internal/helpers"
	"authentication/internal/infra/db"
	models "authentication/internal/models"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"Invalid request body"}`, http.StatusBadRequest)
		return
	}

	var user models.User
	err := db.DB.Get(&user, "SELECT id, name, email, role, password FROM users WHERE email=$1 LIMIT 1", req.Email)
	if err != nil {
		http.Error(w, `{"error":"Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	// Compare password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, `{"error":"Invalid email or password"}`, http.StatusUnauthorized)
		return
	}

	token, err := helpers.GenerateToken(user.ID.String(), user.Role)
	if err != nil {
		http.Error(w, `{"error":"Token generation failed"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"user": map[string]interface{}{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
		"token": token,
	})
}
