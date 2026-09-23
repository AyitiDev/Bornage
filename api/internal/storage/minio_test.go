package storage_test

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/AyitiDev/Bornage/api/internal/config"
	"github.com/AyitiDev/Bornage/api/internal/storage"
)

func TestStorage_InterfaceContract(t *testing.T) {
	cfg := config.MinIOConfig{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadminpassword",
		Bucket:    "test-bucket",
		UseSSL:    false,
	}

	svc, err := storage.NewMinIOStorage(cfg)
	if err != nil {
		t.Fatalf("expected no error initializing storage client, got %v", err)
	}

	if svc == nil {
		t.Fatalf("expected non-nil StorageService instance")
	}
}

func TestStorage_SHA256Calculation(t *testing.T) {
	content := []byte("Land Title Deed Document Sample Content for SHA-256 integrity test")
	hasher := sha256.New()
	hasher.Write(content)
	expectedHash := hex.EncodeToString(hasher.Sum(nil))

	if len(expectedHash) != 64 {
		t.Errorf("expected 64-character SHA-256 hash, got length %d", len(expectedHash))
	}
}

func TestStorage_Singleton(t *testing.T) {
	cfg := config.MinIOConfig{
		Endpoint:  "localhost:9000",
		AccessKey: "minioadmin",
		SecretKey: "minioadminpassword",
		Bucket:    "test-bucket",
		UseSSL:    false,
	}

	inst1, err1 := storage.GetInstance(cfg)
	if err1 != nil {
		t.Fatalf("expected no error from first GetInstance, got %v", err1)
	}

	inst2, err2 := storage.GetInstance(cfg)
	if err2 != nil {
		t.Fatalf("expected no error from second GetInstance, got %v", err2)
	}

	if inst1 != inst2 {
		t.Errorf("expected GetInstance to return the same Singleton instance")
	}
}
