package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/vit6556/avito-internship-assignment/internal/config"
)

const migrationsDir = "file://migrations"

func main() {
	cfg := config.LoadDatabaseConfig()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		cfg.Username, cfg.Password,
		net.JoinHostPort(cfg.Host, cfg.Port),
		cfg.Name,
	)

	m, err := migrate.New(migrationsDir, connString)
	if err != nil {
		log.Fatalf("Error initializing migrations: %v", err)
	}

	flag.Usage = func() {
		log.Println("Usage:")
		log.Println("  migrator up        - Apply all migrations")
		log.Println("  migrator down      - Rollback the last migration")
	}

	flag.Parse()
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	command := flag.Arg(0)
	switch command {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error applying migrations: %v", err)
		}
		log.Println("All migrations applied successfully.")

	case "down":
		if err := m.Down(); err != nil {
			log.Fatalf("Error rolling back migration: %v", err)
		}
		log.Println("Last migration rolled back.")

	default:
		flag.Usage()
		os.Exit(1)
	}
}
