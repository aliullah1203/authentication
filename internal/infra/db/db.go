package db

import (
	"authentication/config"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func ConnectPostgres(DBconfig config.DBconfig) {
	host := DBconfig.Host
	port := DBconfig.Port
	user := DBconfig.User
	password := DBconfig.Password
	dbname := DBconfig.DBName
	sslmode := DBconfig.SSLMode

	if host == "" || port == "" || user == "" || password == "" || dbname == "" {
		log.Fatal("Database environment variables are missing")
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)

	var err error
	DB, err = sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to PostgreSQL: %v", err)
	}

	// Ping to verify live connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping DB: %v", err)
	}

	log.Println("Connected to PostgreSQL successfully")
}
