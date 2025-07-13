package main

import (
	"database/sql"
	"go-rest/internal/database"
	"go-rest/internal/env"
	"log"

	_ "go-rest/docs" // Import generated Swagger docs

	_ "github.com/joho/godotenv/autoload" // Automatically load .env file
	_ "github.com/mattn/go-sqlite3"
	// Swagger handler
	// Swagger files
)

// @title GO Gin Rest API
// @version 1.0
// @description This is a sample Go Gin REST API
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
// @description Enter your JWT Bearer Token in format **bearer &lt;token&gt;**
// @contact.url https://github.com/Mersad-Moghaddam/Complete-GO-REST-API
// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

type Application struct {
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
	app := &Application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "secret"),
		models:    models,
	}

	if err := app.serve(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
