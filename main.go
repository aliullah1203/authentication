package main

import (
	"authentication/config"
	"authentication/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect PostgreSQL
	config.ConnectPostgres()

	r := gin.Default()

	// Register routes
	routes.RegisterRoutes(r)

	r.Run(":8080")
}
