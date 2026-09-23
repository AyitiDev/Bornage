package storage

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/AyitiDev/Bornage/api/internal/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	storageInstance StorageService
	storageOnce     sync.Once
	storageErr      error
)

// MinIOStorage implements StorageService using the MinIO / S3 Go SDK
type MinIOStorage struct {
	client *minio.Client
	bucket string
}

// GetInstance returns a thread safe Singleton instance of the StorageService
func GetInstance(cfg config.MinIOConfig) (StorageService, error) {
	storageOnce.Do(func() {
		storageInstance, storageErr = NewMinIOStorage(cfg)
	})
	return storageInstance, storageErr
}

// NewMinIOStorage creates a new MinIO/S3 storage client and ensures the default bucket exists
func NewMinIOStorage(cfg config.MinIOConfig) (StorageService, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to initialize MinIO client: %w", err)
	}

	storage := &MinIOStorage{
		client: client,
		bucket: cfg.Bucket,
	}

	// Verify or auto create evidence bucket on startup
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := storage.EnsureBucketExists(ctx); err != nil {
		log.Printf("[WARN] Evidence bucket '%s' check returned: %v (will retry on upload)", cfg.Bucket, err)
	} else {
		log.Printf("[INFO] MinIO Storage initialized. Target bucket: '%s'", cfg.Bucket)
	}

	return storage, nil
}

// EnsureBucketExists creates the evidence bucket if it does not already exist
func (s *MinIOStorage) EnsureBucketExists(ctx context.Context) error {
	exists, err := s.client.BucketExists(ctx, s.bucket)
	if err != nil {
		return fmt.Errorf("failed to check if bucket '%s' exists: %w", s.bucket, err)
	}

	if !exists {
		err = s.client.MakeBucket(ctx, s.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket '%s': %w", s.bucket, err)
		}
		log.Printf("[INFO] Created missing bucket: %s", s.bucket)
	}

	return nil
}

// Upload buffers or streams data to compute the SHA-256 integrity hash and uploads the object to MinIO/S3
func (s *MinIOStorage) Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (string, string, error) {
	// Read into memory or calculate hash while reading to guarantee tamper proof audit
	var buf bytes.Buffer
	hasher := sha256.New()
	multiWriter := io.MultiWriter(&buf, hasher)

	written, err := io.Copy(multiWriter, reader)
	if err != nil {
		return "", "", fmt.Errorf("failed to read object data for hashing: %w", err)
	}

	sha256Hash := hex.EncodeToString(hasher.Sum(nil))

	if contentType == "" {
		contentType = "application/octet-stream"
	}

	uploadSize := size
	if uploadSize <= 0 {
		uploadSize = written
	}

	uploadInfo, err := s.client.PutObject(ctx, s.bucket, objectName, &buf, uploadSize, minio.PutObjectOptions{
		ContentType: contentType,
		UserMetadata: map[string]string{
			"sha256": sha256Hash,
		},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to upload object '%s' to bucket '%s': %w", objectName, s.bucket, err)
	}

	log.Printf("[INFO] Stored evidence '%s' (%d bytes, SHA256: %s)", uploadInfo.Key, uploadInfo.Size, sha256Hash)
	return uploadInfo.Key, sha256Hash, nil
}

// GetPresignedURL generates a secure, time limited presigned GET URL for downloading evidence
func (s *MinIOStorage) GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	reqParams := make(url.Values)
	presignedURL, err := s.client.PresignedGetObject(ctx, s.bucket, objectName, expires, reqParams)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL for '%s': %w", objectName, err)
	}
	return presignedURL.String(), nil
}

// Delete permanently removes an evidence object from the bucket
func (s *MinIOStorage) Delete(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object '%s' from bucket '%s': %w", objectName, s.bucket, err)
	}
	return nil
}
