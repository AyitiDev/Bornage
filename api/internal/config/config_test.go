package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/AyitiDev/Bornage/api/internal/config"
)

func TestConfig_LoadDefaults(t *testing.T) {
	// Clear any existing env variables that could interfere
	os.Unsetenv("PORT")
	os.Unsetenv("ENV")
	os.Unsetenv("DB_HOST")
	os.Unsetenv("DB_PORT")
	os.Unsetenv("DB_USER")
	os.Unsetenv("DB_NAME")
	os.Unsetenv("MINIO_BUCKET")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error loading default config, got %v", err)
	}

	if cfg.Server.Port != "8080" {
		t.Errorf("expected default Port '8080', got '%s'", cfg.Server.Port)
	}

	if cfg.Server.Env != "development" {
		t.Errorf("expected default Env 'development', got '%s'", cfg.Server.Env)
	}

	if cfg.Database.Host != "localhost" {
		t.Errorf("expected default DB Host 'localhost', got '%s'", cfg.Database.Host)
	}

	if cfg.Database.Port != "5432" {
		t.Errorf("expected default DB Port '5432', got '%s'", cfg.Database.Port)
	}

	if cfg.Database.Name != "openlandregistry" {
		t.Errorf("expected default DB Name 'openlandregistry', got '%s'", cfg.Database.Name)
	}

	if cfg.MinIO.Bucket != "land-registry-evidence" {
		t.Errorf("expected default MinIO bucket 'land-registry-evidence', got '%s'", cfg.MinIO.Bucket)
	}

	if cfg.LADM.ObjectionPeriodDays != 30 {
		t.Errorf("expected default ObjectionPeriodDays 30, got %d", cfg.LADM.ObjectionPeriodDays)
	}

	if cfg.IsProduction() {
		t.Errorf("expected IsProduction() to be false for development env")
	}
}

func TestConfig_CustomEnv(t *testing.T) {
	t.Setenv("PORT", "9090")
	t.Setenv("ENV", "production")
	t.Setenv("DB_HOST", "db.example.com")
	t.Setenv("DB_PORT", "5433")
	t.Setenv("DB_USER", "customuser")
	t.Setenv("DB_PASSWORD", "secret123")
	t.Setenv("DB_NAME", "customdb")
	t.Setenv("DB_SSLMODE", "require")
	t.Setenv("DB_MAX_OPEN_CONNS", "50")
	t.Setenv("DB_MAX_IDLE_CONNS", "10")
	t.Setenv("DB_CONN_MAX_LIFETIME_MIN", "30")
	t.Setenv("MINIO_ENDPOINT", "s3.example.com")
	t.Setenv("MINIO_BUCKET", "custom-bucket")
	t.Setenv("LADM_OVERLAP_TOLERANCE_METERS", "0.10")
	t.Setenv("LADM_OBJECTION_PERIOD_DAYS", "45")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("expected no error loading custom config, got %v", err)
	}

	if cfg.Server.Port != "9090" {
		t.Errorf("expected Port '9090', got '%s'", cfg.Server.Port)
	}

	if !cfg.IsProduction() {
		t.Errorf("expected IsProduction() to be true for production env")
	}

	if cfg.Database.MaxOpenConns != 50 {
		t.Errorf("expected MaxOpenConns 50, got %d", cfg.Database.MaxOpenConns)
	}

	if cfg.Database.ConnMaxLifetime != 30*time.Minute {
		t.Errorf("expected ConnMaxLifetime 30m, got %v", cfg.Database.ConnMaxLifetime)
	}

	expectedDSN := "host=db.example.com port=5433 user=customuser password=secret123 dbname=customdb sslmode=require"
	if cfg.Database.DSN() != expectedDSN {
		t.Errorf("expected DSN '%s', got '%s'", expectedDSN, cfg.Database.DSN())
	}

	if cfg.LADM.ObjectionPeriodDays != 45 {
		t.Errorf("expected ObjectionPeriodDays 45, got %d", cfg.LADM.ObjectionPeriodDays)
	}
}

func TestConfig_Validation(t *testing.T) {
	tests := []struct {
		name      string
		modifier  func(c *config.Config)
		expectErr bool
	}{
		{
			name: "empty PORT fails validation",
			modifier: func(c *config.Config) {
				c.Server.Port = ""
			},
			expectErr: true,
		},
		{
			name: "empty DB host fails validation",
			modifier: func(c *config.Config) {
				c.Database.Host = ""
			},
			expectErr: true,
		},
		{
			name: "empty DB user fails validation",
			modifier: func(c *config.Config) {
				c.Database.User = ""
			},
			expectErr: true,
		},
		{
			name: "empty DB name fails validation",
			modifier: func(c *config.Config) {
				c.Database.Name = ""
			},
			expectErr: true,
		},
		{
			name: "empty MinIO bucket fails validation",
			modifier: func(c *config.Config) {
				c.MinIO.Bucket = ""
			},
			expectErr: true,
		},
		{
			name: "invalid objection days fails validation",
			modifier: func(c *config.Config) {
				c.LADM.ObjectionPeriodDays = 0
			},
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.Load()
			if err != nil {
				t.Fatalf("unexpected error loading base config: %v", err)
			}
			tt.modifier(cfg)
			valErr := cfg.Validate()
			if tt.expectErr && valErr == nil {
				t.Errorf("expected validation error, but got nil")
			}
			if !tt.expectErr && valErr != nil {
				t.Errorf("unexpected validation error: %v", valErr)
			}
		})
	}
}
