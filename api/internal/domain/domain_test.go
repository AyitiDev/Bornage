package domain_test

import (
	"testing"

	"github.com/AyitiDev/Bornage/api/internal/domain"
	"github.com/google/uuid"
)

func TestNewParty(t *testing.T) {
	t.Run("valid natural person", func(t *testing.T) {
		p, err := domain.NewParty("Jean", "Baptiste", domain.PartyTypeIndividual, "NIF-12345", "ONIExt")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if p.ID == uuid.Nil {
			t.Errorf("expected generated UUID, got nil")
		}
		if p.FullName() != "Jean Baptiste" {
			t.Errorf("expected 'Jean Baptiste', got '%s'", p.FullName())
		}
	})

	t.Run("empty first name fails for individual", func(t *testing.T) {
		_, err := domain.NewParty("", "Dupont", domain.PartyTypeIndividual, "", "")
		if err == nil {
			t.Errorf("expected error for empty first name")
		}
	})

	t.Run("invalid party type fails", func(t *testing.T) {
		_, err := domain.NewParty("Jean", "Baptiste", domain.PartyType("ALIEN"), "", "")
		if err == nil {
			t.Errorf("expected error for invalid party type")
		}
	})
}

func TestNewSpatialUnit(t *testing.T) {
	t.Run("valid spatial unit", func(t *testing.T) {
		geoJSON := `{"type":"Polygon","coordinates":[[[-72.3,18.5],[-72.3,18.6],[-72.2,18.6],[-72.3,18.5]]]}`
		su, err := domain.NewSpatialUnit(geoJSON, 2.5, domain.CaptureMethodSatelliteVisualDelineation)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if su.ID == uuid.Nil {
			t.Errorf("expected generated UUID")
		}
		if su.AccuracyMeters != 2.5 {
			t.Errorf("expected 2.5m accuracy, got %f", su.AccuracyMeters)
		}
	})

	t.Run("empty geometry fails", func(t *testing.T) {
		_, err := domain.NewSpatialUnit("", 2.5, domain.CaptureMethodDroneOrthophoto)
		if err == nil {
			t.Errorf("expected error for empty geometry")
		}
	})

	t.Run("negative accuracy fails", func(t *testing.T) {
		_, err := domain.NewSpatialUnit(`{"type":"Polygon"}`, -1.0, domain.CaptureMethodDroneOrthophoto)
		if err == nil {
			t.Errorf("expected error for negative accuracy")
		}
	})
}

func TestNewRRR(t *testing.T) {
	spatialUnitID := uuid.New()
	partyID := uuid.New()

	t.Run("valid ownership right", func(t *testing.T) {
		rrr, err := domain.NewRRR(spatialUnitID, partyID, domain.RRRTypeOwnership, "1/1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if rrr.ShareFraction != "1/1" {
			t.Errorf("expected 1/1, got %s", rrr.ShareFraction)
		}
	})

	t.Run("nil foreign key fails", func(t *testing.T) {
		_, err := domain.NewRRR(uuid.Nil, partyID, domain.RRRTypeOwnership, "1/1")
		if err == nil {
			t.Errorf("expected error for nil spatial_unit_id")
		}
	})
}

func TestNewClaim(t *testing.T) {
	spatialUnitID := uuid.New()
	claimantID := uuid.New()

	t.Run("generates claim in DRAFT status with tracking code", func(t *testing.T) {
		claim, err := domain.NewClaim(spatialUnitID, claimantID, "Initial title claim")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if claim.Status != domain.ClaimStatusDraft {
			t.Errorf("expected DRAFT status, got %s", claim.Status)
		}
		if len(claim.TrackingCode) < 10 {
			t.Errorf("expected valid tracking code format, got '%s'", claim.TrackingCode)
		}
	})
}

func TestNewSource(t *testing.T) {
	claimID := uuid.New()
	validHash := "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	t.Run("valid source document", func(t *testing.T) {
		s, err := domain.NewSource(claimID, "claims/123/deed.pdf", validHash, "application/pdf", "deed.pdf", 1024, "Old title deed")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if s.SHA256Hash != validHash {
			t.Errorf("expected sha256 hash match")
		}
	})

	t.Run("invalid hash length fails", func(t *testing.T) {
		_, err := domain.NewSource(claimID, "claims/123/deed.pdf", "short-hash", "application/pdf", "deed.pdf", 1024, "")
		if err == nil {
			t.Errorf("expected error for invalid sha256 length")
		}
	})
}

func TestNewWitness(t *testing.T) {
	claimID := uuid.New()

	t.Run("valid witness", func(t *testing.T) {
		w, err := domain.NewWitness(claimID, "Marc Antoine", "NIF-8888", "sigs/w1.png", "photos/w1.jpg", "Confirms eastern boundary")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if w.FullName != "Marc Antoine" {
			t.Errorf("expected 'Marc Antoine', got '%s'", w.FullName)
		}
	})

	t.Run("empty full name fails", func(t *testing.T) {
		_, err := domain.NewWitness(claimID, "", "", "", "", "")
		if err == nil {
			t.Errorf("expected error for empty witness full name")
		}
	})
}

func TestNewAuditLog(t *testing.T) {
	t.Run("valid audit log", func(t *testing.T) {
		changes := map[string]interface{}{"status": "PUBLISHED_FOR_OBJECTION"}
		log, err := domain.NewAuditLog("actor-123", "UPDATE_STATUS", "Claim", "claim-456", changes, "192.168.1.1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if log.Action != "UPDATE_STATUS" {
			t.Errorf("expected 'UPDATE_STATUS', got '%s'", log.Action)
		}
		if len(log.ChangesPayload) == 0 {
			t.Errorf("expected JSON changes payload")
		}
	})
}
