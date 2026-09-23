package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Source represents documentary or photographic evidence attached to a claim (LADM LA_AdministrativeSource)
// The sha256_hash guarantees tamper detection for uploaded legal evidence
type Source struct {
	ID               uuid.UUID `json:"id" db:"id"`
	ClaimID          uuid.UUID `json:"claim_id" db:"claim_id"`
	S3Key            string    `json:"s3_key" db:"s3_key"`
	FileType         string    `json:"file_type" db:"file_type"`
	OriginalFilename string    `json:"original_filename,omitempty" db:"original_filename"`
	SHA256Hash       string    `json:"sha256_hash" db:"sha256_hash"`
	FileSizeBytes    int64     `json:"file_size_bytes" db:"file_size_bytes"`
	Description      string    `json:"description,omitempty" db:"description"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// NewSource creates and validates a new Source document record (Factory Method pattern)
func NewSource(claimID uuid.UUID, s3Key, sha256Hash, fileType, originalFilename string, sizeBytes int64, description string) (*Source, error) {
	if claimID == uuid.Nil {
		return nil, errors.New("claim_id cannot be nil")
	}
	s3Key = strings.TrimSpace(s3Key)
	if s3Key == "" {
		return nil, errors.New("s3_key cannot be empty")
	}
	sha256Hash = strings.TrimSpace(sha256Hash)
	if len(sha256Hash) != 64 {
		return nil, errors.New("sha256_hash must be a valid 64-character hex string")
	}
	fileType = strings.TrimSpace(fileType)
	if fileType == "" {
		return nil, errors.New("file_type cannot be empty")
	}

	return &Source{
		ID:               uuid.New(),
		ClaimID:          claimID,
		S3Key:            s3Key,
		FileType:         fileType,
		OriginalFilename: strings.TrimSpace(originalFilename),
		SHA256Hash:       sha256Hash,
		FileSizeBytes:    sizeBytes,
		Description:      strings.TrimSpace(description),
		CreatedAt:        time.Now().UTC(),
	}, nil
}
