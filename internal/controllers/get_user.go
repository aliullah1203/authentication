package controllers

import (
	"authentication/config"
	"authentication/internal/helpers"
	models "authentication/internal/models"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Get single user by ID (ADMIN or same user)
func GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	requestedUserID := vars["id"]

	claims, ok := r.Context().Value("claims").(*helpers.Claims)
	if !ok {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	if claims.Role != "ADMIN" && claims.UserID != requestedUserID {
		http.Error(w, `{"error":"Unauthorized"}`, http.StatusUnauthorized)
		return
	}

	var user models.User
	err := config.DB.Get(&user, `SELECT * FROM users WHERE id=$1`, requestedUserID)
	if err != nil {
		http.Error(w, `{"error":"User not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
