package controllers

import (
	"encoding/json"
	"net/http"
)

// Logout simply informs the client to delete the token
func Logout(w http.ResponseWriter, r *http.Request) {
	// you can clear cookies if storing JWT in cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logged out successfully",
	})
}
