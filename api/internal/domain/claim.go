package domain

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

// Claim represents a land title registration application in the state machine lifecycle
// Mandatory Rule: Must pass through PUBLISHED_FOR_OBJECTION before reaching REGISTERED
type Claim struct {
	ID               uuid.UUID   `json:"id" db:"id"`
	TrackingCode     string      `json:"tracking_code" db:"tracking_code"`
	Status           ClaimStatus `json:"status" db:"status"`
	SpatialUnitID    uuid.UUID   `json:"spatial_unit_id" db:"spatial_unit_id"`
	ClaimantID       uuid.UUID   `json:"claimant_id" db:"claimant_id"`
	ObjectionEndDate *time.Time  `json:"objection_end_date,omitempty" db:"objection_end_date"`
	Notes            string      `json:"notes,omitempty" db:"notes"`
	CreatedAt        time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at" db:"updated_at"`
}

// NewClaim creates and initializes a new Claim in DRAFT status with a unique tracking code (Factory Method pattern)
func NewClaim(spatialUnitID, claimantID uuid.UUID, notes string) (*Claim, error) {
	if spatialUnitID == uuid.Nil {
		return nil, errors.New("spatial_unit_id cannot be nil")
	}
	if claimantID == uuid.Nil {
		return nil, errors.New("claimant_id cannot be nil")
	}

	trackingCode, err := generateTrackingCode()
	if err != nil {
		return nil, fmt.Errorf("failed to generate tracking code: %w", err)
	}

	now := time.Now().UTC()
	return &Claim{
		ID:            uuid.New(),
		TrackingCode:  trackingCode,
		Status:        ClaimStatusDraft,
		SpatialUnitID: spatialUnitID,
		ClaimantID:    claimantID,
		Notes:         notes,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// generateTrackingCode generates a human readable unique reference: CLM-YYYY-<6 RANDOM ALPHANUMERIC>
func generateTrackingCode() (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	year := time.Now().Year()
	b := make([]byte, 6)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return fmt.Sprintf("CLM-%d-%s", year, string(b)), nil
}
