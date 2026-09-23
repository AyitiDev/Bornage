package domain

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Party represents a person, organization, or state holding land rights (LADM LA_Party)
// It maintains external identity references without duplicating sensitive civil registry data
type Party struct {
	ID             uuid.UUID `json:"id" db:"id"`
	FirstName      string    `json:"first_name" db:"first_name"`
	LastName       string    `json:"last_name" db:"last_name"`
	PartyType      PartyType `json:"party_type" db:"party_type"`
	ExternalID     string    `json:"external_id,omitempty" db:"external_id"`
	ExternalSystem string    `json:"external_system,omitempty" db:"external_system"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// NewParty creates and validates a new Party entity (Factory Method pattern)
func NewParty(firstName, lastName string, partyType PartyType, externalID, externalSystem string) (*Party, error) {
	firstName = strings.TrimSpace(firstName)
	lastName = strings.TrimSpace(lastName)

	if firstName == "" && partyType != PartyTypeState {
		return nil, errors.New("first_name cannot be empty")
	}

	if !partyType.IsValid() {
		return nil, errors.New("invalid party_type")
	}

	now := time.Now().UTC()
	return &Party{
		ID:             uuid.New(),
		FirstName:      firstName,
		LastName:       lastName,
		PartyType:      partyType,
		ExternalID:     strings.TrimSpace(externalID),
		ExternalSystem: strings.TrimSpace(externalSystem),
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}

// FullName returns the concatenated party name
func (p *Party) FullName() string {
	if p.LastName == "" {
		return p.FirstName
	}
	return p.FirstName + " " + p.LastName
}
