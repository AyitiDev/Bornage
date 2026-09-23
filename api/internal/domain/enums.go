package domain

import "fmt"

// ClaimStatus represents the lifecycle state machine of a land registration claim
type ClaimStatus string

const (
	ClaimStatusDraft                 ClaimStatus = "DRAFT"
	ClaimStatusSubmitted             ClaimStatus = "SUBMITTED"
	ClaimStatusUnderReview           ClaimStatus = "UNDER_REVIEW"
	ClaimStatusPublishedForObjection ClaimStatus = "PUBLISHED_FOR_OBJECTION"
	ClaimStatusRegistered            ClaimStatus = "REGISTERED"
	ClaimStatusReturnedToField       ClaimStatus = "RETURNED_TO_FIELD"
	ClaimStatusRejected              ClaimStatus = "REJECTED"
	ClaimStatusDisputed              ClaimStatus = "DISPUTED"
)

func (s ClaimStatus) IsValid() bool {
	switch s {
	case ClaimStatusDraft, ClaimStatusSubmitted, ClaimStatusUnderReview,
		ClaimStatusPublishedForObjection, ClaimStatusRegistered,
		ClaimStatusReturnedToField, ClaimStatusRejected, ClaimStatusDisputed:
		return true
	default:
		return false
	}
}

// RRRType represents the Rights, Restrictions, and Responsibilities (LADM LA_RRR)
type RRRType string

const (
	RRRTypeOwnership          RRRType = "OWNERSHIP"
	RRRTypeLease              RRRType = "LEASE"
	RRRTypeUsufruct           RRRType = "USUFRUCT"
	RRRTypeInformalOccupation RRRType = "INFORMAL_OCCUPATION"
	RRRTypeEasement           RRRType = "EASEMENT"
	RRRTypeMortgage           RRRType = "MORTGAGE"
	RRRTypeRestriction        RRRType = "RESTRICTION"
)

func (r RRRType) IsValid() bool {
	switch r {
	case RRRTypeOwnership, RRRTypeLease, RRRTypeUsufruct,
		RRRTypeInformalOccupation, RRRTypeEasement, RRRTypeMortgage, RRRTypeRestriction:
		return true
	default:
		return false
	}
}

// PartyType represents the legal nature of a tenure holder (LADM LA_Party)
type PartyType string

const (
	PartyTypeIndividual   PartyType = "INDIVIDUAL"
	PartyTypeOrganization PartyType = "ORGANIZATION"
	PartyTypeState        PartyType = "STATE"
)

func (p PartyType) IsValid() bool {
	switch p {
	case PartyTypeIndividual, PartyTypeOrganization, PartyTypeState:
		return true
	default:
		return false
	}
}

// CaptureMethod represents the technical survey method used to delineate boundaries
type CaptureMethod string

const (
	CaptureMethodSatelliteVisualDelineation CaptureMethod = "SATELLITE_VISUAL_DELINEATION"
	CaptureMethodGPSWalkedPerimeter         CaptureMethod = "GPS_WALKED_PERIMETER"
	CaptureMethodSurveyTotalStation         CaptureMethod = "SURVEY_TOTAL_STATION"
	CaptureMethodDroneOrthophoto            CaptureMethod = "DRONE_ORTHOPHOTO"
)

func (c CaptureMethod) IsValid() bool {
	switch c {
	case CaptureMethodSatelliteVisualDelineation, CaptureMethodGPSWalkedPerimeter,
		CaptureMethodSurveyTotalStation, CaptureMethodDroneOrthophoto:
		return true
	default:
		return false
	}
}

// ValidateEnum checks if a generic string matches expected enum domain values
func ValidateEnum[T interface{ IsValid() bool }](val T, typeName string) error {
	if !val.IsValid() {
		return fmt.Errorf("invalid %s value", typeName)
	}
	return nil
}
