package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/AyitiDev/Bornage/api/internal/config"
	_ "github.com/lib/pq" // PostgreSQL driver
)

var (
	instance *sql.DB
	once     sync.Once
	initErr  error
)

// GetInstance implements the Singleton Pattern with sync.Once to provide a thread-safe,
// lazily initialized shared connection pool across the application lifecycle
func GetInstance(cfg config.DatabaseConfig) (*sql.DB, error) {
	once.Do(func() {
		instance, initErr = NewPostgresDB(cfg)
	})
	return instance, initErr
}

// NewPostgresDB creates and configures a new PostgreSQL/PostGIS connection pool
// It sets pool limits, verifies connectivity with a ping, and validates PostGIS extension presence
func NewPostgresDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Configure connection pool parameters for high throughput and resource management
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Verify connectivity with a timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database at %s:%s: %w", cfg.Host, cfg.Port, err)
	}

	// Verify PostGIS extension is available
	var postgisVersion string
	query := "SELECT PostGIS_Version();"
	if err := db.QueryRowContext(ctx, query).Scan(&postgisVersion); err != nil {
		log.Printf("[WARN] PostGIS extension check returned: %v (ensure PostGIS is enabled on db %s)", err, cfg.Name)
	} else {
		log.Printf("[INFO] Database connected successfully. PostGIS version: %s", postgisVersion)
	}

	return db, nil
}

// Close gracefully terminates the singleton database connection pool
func Close() error {
	if instance != nil {
		log.Println("[INFO] Closing database connection pool...")
		return instance.Close()
	}
	return nil
}
