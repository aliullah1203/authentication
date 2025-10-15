package main

import (
	"authentication/config"
	"authentication/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file (no fatal if not present, but recommended)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, rely on env variables")
	}

	// Connect PostgreSQL
	config.ConnectPostgres()

	// Create Gin router
	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r)

	// Port fallback
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)
}
