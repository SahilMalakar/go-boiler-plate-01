package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/sahil_malakar/production_grade_golang_setup/internal/config"
)

func main() {
	cfg := config.MustLoad()

	if len(os.Args) < 2 {
		log.Fatal("usage: migrate <up | down>")
	}

	// Creates the migration instance using the local migrations directory
	// and the configured PostgreSQL database.
	m, err := migrate.New(
		"file://migrations",
		cfg.DatabaseURL,
	)
	if err != nil {
		log.Fatalf("migration.new: %v", err)
	}
	defer m.Close()

	fmt.Printf("Running migration %s \n", os.Args[1])

	switch os.Args[1] {
	case "up":
		if err := m.Up(); err != nil {
			log.Fatalf("migration.up: %v", err)
		}

	case "down":
		if err := m.Steps(-1); err != nil {
			log.Fatalf("migration.down: %v", err)
		}

	default:
		log.Fatalf("unknown command: %s", os.Args[1])
	}

	fmt.Println("Migration completed")
}
