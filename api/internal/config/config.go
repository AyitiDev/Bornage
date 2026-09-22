package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all strongly typed configuration for the Open Land Registry backend
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	MinIO    MinIOConfig
	LADM     LADMConfig
}

// ServerConfig contains HTTP server runtime settings
type ServerConfig struct {
	Port         string
	Env          string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig contains PostgreSQL / PostGIS connection parameters
type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// DSN returns the PostgreSQL connection string
func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode,
	)
}

// MinIOConfig contains S3/MinIO object storage settings for document and evidence storage
type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

// LADMConfig contains land registry domain business rules and tolerance thresholds
type LADMConfig struct {
	OverlapToleranceMeters float64 // Fit For Purpose spatial boundary tolerance threshold (e.x. 0.05m or 3m)
	ObjectionPeriodDays    int     // Mandatory minimum days for public objection notice period
}

// Load reads configuration from environment variables, applies defaults, and validates fields
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			Env:          getEnv("ENV", "development"),
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 10 * time.Second,
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgrespassword"),
			Name:            getEnv("DB_NAME", "openlandregistry"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: time.Duration(getEnvAsInt("DB_CONN_MAX_LIFETIME_MIN", 15)) * time.Minute,
		},
		MinIO: MinIOConfig{
			Endpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
			AccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("MINIO_SECRET_KEY", "minioadminpassword"),
			Bucket:    getEnv("MINIO_BUCKET", "land-registry-evidence"),
			UseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
		},
		LADM: LADMConfig{
			OverlapToleranceMeters: getEnvAsFloat("LADM_OVERLAP_TOLERANCE_METERS", 0.05),
			ObjectionPeriodDays:    getEnvAsInt("LADM_OBJECTION_PERIOD_DAYS", 30),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation error: %w", err)
	}

	return cfg, nil
}

// Validate checks that all required configuration properties are valid
func (c *Config) Validate() error {
	if c.Server.Port == "" {
		return fmt.Errorf("PORT cannot be empty")
	}
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST cannot be empty")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER cannot be empty")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME cannot be empty")
	}
	if c.MinIO.Bucket == "" {
		return fmt.Errorf("MINIO_BUCKET cannot be empty")
	}
	if c.LADM.ObjectionPeriodDays < 1 {
		return fmt.Errorf("LADM_OBJECTION_PERIOD_DAYS must be at least 1 day")
	}
	return nil
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsFloat(key string, defaultVal float64) float64 {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseFloat(valStr, 64); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}
