package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/AyitiDev/Bornage/api/internal/db"
	_ "github.com/lib/pq"
)

func main() {
	var command string
	var migrationsDir string

	flag.StringVar(&command, "cmd", "up", "Migration command to execute: 'up' or 'down'")
	flag.StringVar(&migrationsDir, "dir", "./migrations", "Path to migrations directory")
	flag.Parse()

	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPass := getEnv("DB_PASSWORD", "postgrespassword")
	dbName := getEnv("DB_NAME", "openlandregistry")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)

	sqlDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer sqlDB.Close()

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Database ping failed: %v", err)
	}

	log.Printf("Connected to PostgreSQL database '%s' at %s:%s", dbName, dbHost, dbPort)

	migrator := db.NewMigrator(sqlDB, migrationsDir)
	if err := migrator.Run(command); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Printf("Migration command '%s' completed successfully.", command)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
