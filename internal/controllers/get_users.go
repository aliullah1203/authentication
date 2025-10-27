package controllers

import (
	"authentication/internal/helpers"
	"authentication/internal/infra/db"
	models "authentication/internal/models"
	"encoding/json"
	"net/http"
)

// Get all users (ADMIN/SUPER_ADMIN)
func GetUsers(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value("claims").(*helpers.Claims)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if claims.Role != "ADMIN" && claims.Role != "SUPER_ADMIN" {
		http.Error(w, `{"error":"Forbidden"}`, http.StatusForbidden)
		return
	}

	var users []models.User
	err := db.DB.Select(&users, `SELECT * FROM users`)
	if err != nil {
		http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"users": users,
	})
}
