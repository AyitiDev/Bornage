package storage

import (
	"context"
	"io"
	"time"
)

// StorageService defines the contract for document and evidence storage
// This interface complies with the Dependency Inversion Principle (DIP)
type StorageService interface {
	// Upload stores an object and returns its s3_key, calculated sha256_hash, and error
	Upload(ctx context.Context, objectName string, reader io.Reader, size int64, contentType string) (s3Key string, sha256Hash string, err error)

	// GetPresignedURL generates a time-limited download URL for a stored evidence file
	GetPresignedURL(ctx context.Context, objectName string, expires time.Duration) (string, error)

	// Delete removes an object from the bucket
	Delete(ctx context.Context, objectName string) error

	// EnsureBucketExists checks if the evidence bucket exists, creating it if necessary
	EnsureBucketExists(ctx context.Context) error
}
