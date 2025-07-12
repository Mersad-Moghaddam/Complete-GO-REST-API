package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Please provide a migration direction: `up` or `down`")
	}

	direction := os.Args[1]

	// Connect to SQLite database
	db, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatalf("Failed to enable foreign keys: %v", err)
	}

	defer db.Close()

	// Set up migrate with sqlite3
	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatalf("Failed to create SQLite migration driver: %v", err)
	}

	sourceDriver, err := (&file.File{}).Open("cmd/migrate/migrations")
	if err != nil {
		log.Fatalf("Failed to open migration source: %v", err)
	}

	m, err := migrate.NewWithInstance(
		"file",
		sourceDriver,
		"sqlite3",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	switch direction {
	case "up":
		err := m.Up()
		if err != nil {
			if err.Error() == "dirty database" || err.Error() == "Dirty database version 3. Fix and force version." {
				version, dirty, verr := m.Version()
				if verr != nil {
					log.Fatalf("Failed to get current migration version: %v", verr)
				}
				if dirty {
					log.Printf("Dirty database detected at version %d. Forcing clean state...", version)
					if forceErr := m.Force(int(version)); forceErr != nil {
						log.Fatalf("Failed to force migration version: %v", forceErr)
					}
					log.Println("Forced migration version. Retrying migration...")
					err = m.Up()
					if err != nil && err != migrate.ErrNoChange {
						log.Fatalf("Migration failed again after forcing: %v", err)
					}
				}
			} else if err != migrate.ErrNoChange {
				log.Fatalf("Failed to apply migrations up: %v", err)
			}
		}
		log.Println("Migrations applied successfully (up)")
	case "down":
		err := m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations down: %v", err)
		}
		log.Println("Migrations applied successfully (down)")
	default:
		log.Fatalf("Invalid direction: %s. Use `up` or `down`", direction)
	}
}
