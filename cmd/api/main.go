package main

import (
	"database/sql"
	"go-rest/internal/database"
	"go-rest/internal/env"
	"log"

	_ "github.com/joho/godotenv/autoload" // Automatically load .env file
	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	port      int
	jwtSecret string
	models    database.Models
}

func main() {
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}
	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "secret"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
