package cmd

import (
	"authentication/config"
	"authentication/internal/infra/db"
	"authentication/internal/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func Service() {
	// Load .env file
	config.LoadEnv()

	// Connect to postgres
	cfg := config.DBconfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}
	db.ConnectPostgres(cfg)

	// Initialize Google OAuth config from env
	if err := services.InitGoogleOAuthFromEnv(); err != nil {
		log.Fatalf("Google OAuth config error: %v", err)
	}

	// Create net/http mux and register routes
	router := mux.NewRouter()
	RegisterHTTPRoutes(router)

	// Port fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server run on:%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal(err)
	}
}
