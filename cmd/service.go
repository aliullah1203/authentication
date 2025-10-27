package cmd

import (
	"authentication/config"
	"authentication/internal/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func Service() {
	// Load .env file (no fatal if not present, but recommended)
	config.LoadEnv()

	// Initialize JWT secret and database
	config.ConnectPostgres()

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
