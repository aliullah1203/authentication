package main

import (
    "authentication/config"
    "authentication/routes"
    "authentication/services"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
)

func main() {
    // Load .env file (no fatal if not present, but recommended)
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, rely on env variables")
    }

    // Initialize JWT secret and database
    config.ConnectPostgres()

    // Initialize Google OAuth config from env
    if err := services.InitGoogleOAuthFromEnv(); err != nil {
        log.Fatalf("Google OAuth config error: %v", err)
    }

    // Create net/http mux and register routes
    mux := http.NewServeMux()
    routes.RegisterHTTPRoutes(mux)

    // Port fallback
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    log.Printf("Listening on :%s", port)
    if err := http.ListenAndServe(":"+port, mux); err != nil {
        log.Fatal(err)
    }
}
