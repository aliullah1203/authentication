package services

import (
	"authentication/config"
	"authentication/helpers"
	"authentication/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Google OAuth2 config
var googleOauthConfig = &oauth2.Config{
	ClientID:     "YOUR_GOOGLE_CLIENT_ID",
	ClientSecret: "YOUR_GOOGLE_CLIENT_SECRET",
	RedirectURL:  "http://localhost:8080/api/oauth/google/callback",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

// GenerateAuthURL generates URL to redirect user to Google login
func GetGoogleLoginURL(state string) string {
	return googleOauthConfig.AuthCodeURL(state)
}

// HandleGoogleCallback handles the callback from Google
func HandleGoogleCallback(code string) (*models.User, error) {
	// Exchange code for token
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %v", err)
	}

	// Get user info from Google API
	client := googleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, fmt.Errorf("failed to decode user info: %v", err)
	}

	// Check if user already exists in DB
	var user models.User
	err = config.DB.Get(&user, "SELECT * FROM users WHERE email=$1", googleUser.Email)
	if err != nil {
		// User does not exist, create new one
		user = models.User{
			ID:                 uuid.New(),
			Name:               googleUser.Name,
			Email:              googleUser.Email,
			Role:               "CUSTOMER",
			Status:             "ACTIVE",
			SubscriptionStatus: "SUBSCRIBED",
			CreatedAt:          time.Now(),
			UpdatedAt:          time.Now(),
		}
		_, err := config.DB.NamedExec(`INSERT INTO users (id, name, email, role, status, subscription_status, created_at, updated_at)
			VALUES (:id, :name, :email, :role, :status, :subscription_status, :created_at, :updated_at)`, &user)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %v", err)
		}
	}

	// Generate JWT tokens
	accessToken, refreshToken, _ := helpers.GenerateTokens(user.ID.String(), user.Role)

	fmt.Println("AccessToken:", accessToken)
	fmt.Println("RefreshToken:", refreshToken)

	return &user, nil
}
