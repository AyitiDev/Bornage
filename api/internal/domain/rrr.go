package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// RRR represents Rights, Restrictions, and Responsibilities over land (LADM LA_RRR)
// Links a Party to a SpatialUnit with a specific tenure type and ownership share
type RRR struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	SpatialUnitID uuid.UUID  `json:"spatial_unit_id" db:"spatial_unit_id"`
	PartyID       uuid.UUID  `json:"party_id" db:"party_id"`
	RRRType       RRRType    `json:"rrr_type" db:"rrr_type"`
	ShareFraction string     `json:"share_fraction" db:"share_fraction"`
	ValidFrom     time.Time  `json:"valid_from" db:"valid_from"`
	ValidTo       *time.Time `json:"valid_to,omitempty" db:"valid_to"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// NewRRR creates and validates a new Rights/Restrictions/Responsibilities record (Factory Method pattern)
func NewRRR(spatialUnitID, partyID uuid.UUID, rrrType RRRType, shareFraction string) (*RRR, error) {
	if spatialUnitID == uuid.Nil {
		return nil, errors.New("spatial_unit_id cannot be nil")
	}
	if partyID == uuid.Nil {
		return nil, errors.New("party_id cannot be nil")
	}
	if !rrrType.IsValid() {
		return nil, errors.New("invalid rrr_type")
	}

	shareFraction = strings.TrimSpace(shareFraction)
	if shareFraction == "" {
		shareFraction = "1/1" // Default full share
	}

	now := time.Now().UTC()
	return &RRR{
		ID:            uuid.New(),
		SpatialUnitID: spatialUnitID,
		PartyID:       partyID,
		RRRType:       rrrType,
		ShareFraction: shareFraction,
		ValidFrom:     now,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}
