package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Witness represents neighboring occupants who testify and confirm parcel boundaries on site
type Witness struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ClaimID        uuid.UUID `json:"claim_id" db:"claim_id"`
	FullName       string    `json:"full_name" db:"full_name"`
	NationalID     string    `json:"national_id,omitempty" db:"national_id"`
	SignatureS3Key string    `json:"signature_s3_key,omitempty" db:"signature_s3_key"`
	PhotoS3Key     string    `json:"photo_s3_key,omitempty" db:"photo_s3_key"`
	Comments       string    `json:"comments,omitempty" db:"comments"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// NewWitness creates and validates a new boundary Witness record (Factory Method pattern)
func NewWitness(claimID uuid.UUID, fullName, nationalID, signatureS3Key, photoS3Key, comments string) (*Witness, error) {
	if claimID == uuid.Nil {
		return nil, errors.New("claim_id cannot be nil")
	}
	fullName = strings.TrimSpace(fullName)
	if fullName == "" {
		return nil, errors.New("witness full_name cannot be empty")
	}

	return &Witness{
		ID:             uuid.New(),
		ClaimID:        claimID,
		FullName:       fullName,
		NationalID:     strings.TrimSpace(nationalID),
		SignatureS3Key: strings.TrimSpace(signatureS3Key),
		PhotoS3Key:     strings.TrimSpace(photoS3Key),
		Comments:       strings.TrimSpace(comments),
		CreatedAt:      time.Now().UTC(),
	}, nil
}
