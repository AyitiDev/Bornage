package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// Migrator handles database schema migrations
type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrator creates a new instance of Migrator
func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// Run executes migration files matching the specified direction (up or down)
func (m *Migrator) Run(direction string) error {
	pattern := filepath.Join(m.migrationsDir, fmt.Sprintf("*.%s.sql", direction))
	files, err := filepath.Glob(pattern)
	if err != nil || len(files) == 0 {
		return fmt.Errorf("no migration files found for pattern: %s", pattern)
	}

	for _, file := range files {
		log.Printf("Executing migration: %s", filepath.Base(file))
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", file, err)
		}

		if _, err := tx.Exec(string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit transaction for %s: %w", file, err)
		}

		log.Printf("Successfully applied migration: %s", filepath.Base(file))
	}

	return nil
}
